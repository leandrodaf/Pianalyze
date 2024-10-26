package cmd

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"

	channelMIDI "github.com/leandrodaf/midi-client/internal/channel"
	"github.com/leandrodaf/midi-client/internal/contracts/channel"
	"github.com/leandrodaf/midi-client/internal/contracts/logger"
	"github.com/leandrodaf/midi-client/internal/contracts/midi"
	"github.com/leandrodaf/midi-client/internal/listeners"
	"github.com/leandrodaf/midi-client/pkg/environment"
	zaplogger "github.com/leandrodaf/midi-client/pkg/logger"
	midiClient "github.com/leandrodaf/midi-client/pkg/midi"
	"github.com/leandrodaf/midi-client/pkg/pubsub"
)

// SetupLogger initializes the logging system.
func SetupLogger() (logger.Logger, error) {
	return zaplogger.New(logger.Options{
		Environment: environment.IsProduction(),
	})
}

// SetupMIDIClient configures the MIDI client.
func SetupMIDIClient(logger logger.Logger) (midi.ClientMIDI, error) {
	return midiClient.NewMIDIClient(logger)
}

// SetupEventChannel creates the MIDI event channel.
func SetupEventChannel(bufferSize int) channel.MIDIEventPublisher {
	return channelMIDI.NewEventChannel(bufferSize)
}

// SetupHub configures and returns the Hub and topic names.
func SetupHub() (*pubsub.Hub, string, string) {
	hub := pubsub.NewHub()
	topicName := "midi_events"
	processedTopicName := "processed_events"
	return hub, topicName, processedTopicName
}

// SetupListenerManager configures and returns the ListenerManager with registered listeners.
func SetupListenerManager(
	hub *pubsub.Hub,
	topicName string,
	processedTopicName string,
	logger logger.Logger,
	latencyInterval time.Duration,
) *listeners.ListenerManager {
	listenerManager := listeners.NewListenerManager(hub, topicName, logger)

	// Create and register listeners, all publishing to the same output topic.
	chordListener := &listeners.ChordDetectionListener{OutputTopic: processedTopicName, Hub: hub}
	velocityListener := &listeners.VelocityAnalyzerListener{OutputTopic: processedTopicName, Hub: hub}
	listenerManager.Register(chordListener)
	listenerManager.Register(velocityListener)

	go startLatencyTicker(logger, listenerManager, latencyInterval)

	return listenerManager
}

// startLatencyTicker periodically calculates and logs the average latency.
func startLatencyTicker(logger logger.Logger, listenerManager *listeners.ListenerManager, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		latency := listenerManager.AverageLatency()
		latencyRounded := math.Round(latency*1000) / 1000

		logger.Info("Calculating average latency",
			logger.Field().Float64("latency_average_ns", latencyRounded))
	}
}

// SetupDevice selects and configures the MIDI device.
func SetupDevice(adapter midi.ClientMIDI) (int, error) {
	devices, err := adapter.ListDevices()
	if err != nil {
		return 0, err
	}
	fmt.Println("Available MIDI devices:")
	for i, device := range devices {
		fmt.Printf("[%d] %s\n", i, device)
	}
	var deviceID int
	fmt.Print("Choose a MIDI device: ")
	_, err = fmt.Scanf("%d", &deviceID)
	if err != nil {
		return 0, err
	}
	err = adapter.SelectDevice(deviceID)
	if err != nil {
		return deviceID, err
	}
	return deviceID, nil
}

// StartEventPublishing begins capturing and publishing MIDI events.
func StartEventPublishing(eventChannel channel.MIDIEventPublisher, hub *pubsub.Hub, topicName string) {
	go func() {
		for event := range eventChannel.Events() {
			// Capture the current time in nanoseconds
			eventTimestamp := time.Now().UnixNano()

			// Prepare the data buffer: 8 bytes for the timestamp and 3 bytes for the MIDI data
			data := make([]byte, 8+3)
			binary.BigEndian.PutUint64(data[:8], uint64(eventTimestamp))
			data[8] = byte(event.Command)
			data[9] = byte(event.Note)
			data[10] = byte(event.Velocity)

			// Create the event message and publish to the hub
			msg := pubsub.Message{Data: data}
			hub.Publish(topicName, msg)
		}
	}()
}
