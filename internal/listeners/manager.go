package listeners

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/leandrodaf/midi-client/internal/contracts/logger"
	"github.com/leandrodaf/midi-client/internal/entity"
	"github.com/leandrodaf/midi-client/pkg/pubsub"
)

// ListenerManager manages a set of listeners and processes events for them.
type ListenerManager struct {
	hub          *pubsub.Hub
	topicName    string
	listeners    []Listener
	logger       logger.Logger
	totalLatency int64
	messageCount int64
}

// NewListenerManager creates a new instance of ListenerManager.
func NewListenerManager(hub *pubsub.Hub, topicName string, logger logger.Logger) *ListenerManager {
	return &ListenerManager{
		hub:       hub,
		topicName: topicName,
		logger:    logger,
	}
}

// Register adds a listener to the manager.
func (lm *ListenerManager) Register(listener Listener) {
	lm.listeners = append(lm.listeners, listener)
}

// Start begins listening and processing events for all registered listeners.
func (lm *ListenerManager) Start() {
	eventChannel := lm.hub.Subscribe(lm.topicName)
	go func() {
		defer lm.hub.Unsubscribe(lm.topicName, eventChannel)
		for msg := range eventChannel {
			lm.processMessage(msg)
		}
	}()
}

// processMessage distributes the MIDI event to all registered listeners.
func (lm *ListenerManager) processMessage(msg pubsub.Message) {
	// Decode the MIDI event from the message data
	midiEvent, err := decodeMIDIEvent(msg.Data)
	if err != nil {
		lm.logger.Error("Failed to decode MIDI event", lm.logger.Field().Error("decode_error", err))
		return
	}

	// Log the MIDI event details only if necessary for debugging

	lm.logger.Debug("Processing MIDI event",
		lm.logger.Field().Uint64("timestamp", midiEvent.Timestamp),
		lm.logger.Field().Uint8("command", midiEvent.Command),
		lm.logger.Field().Uint8("note", midiEvent.Note),
		lm.logger.Field().Uint8("velocity", midiEvent.Velocity),
	)

	// Measure and log latency for monitoring purposes
	measureEventLatency(msg, lm)

	// Distribute the message to all listeners
	for _, listener := range lm.listeners {
		listener.ProcessMessage(msg)
	}
}

// measureEventLatency calculates and logs the latency of an event.
func measureEventLatency(msg pubsub.Message, lm *ListenerManager) {
	eventTimestamp := int64(binary.BigEndian.Uint64(msg.Data[:8]))
	receiveTimestamp := time.Now().UnixNano()
	latency := receiveTimestamp - eventTimestamp

	lm.totalLatency += latency
	lm.messageCount++

	// Log latency at the info level if it is unusually high, otherwise at the debug level
	if latency > 1e6 { // Example threshold: 1 millisecond
		lm.logger.Info("High listener latency detected", lm.logger.Field().Int64("latency_ns", latency))
	} else {
		lm.logger.Debug("Listener latency", lm.logger.Field().Int64("latency_ns", latency))
	}
}

// decodeMIDIEvent decodes a MIDI event from raw data.
func decodeMIDIEvent(data []byte) (entity.MIDI, error) {
	if len(data) < 11 { // Adjust the length based on the actual data format
		return entity.MIDI{}, fmt.Errorf("data too short to decode MIDI event")
	}

	// Assuming the first 8 bytes are the timestamp
	midiEvent := entity.MIDI{
		Timestamp: binary.BigEndian.Uint64(data[:8]),
		Command:   data[8],
		Note:      data[9],
		Velocity:  data[10],
	}

	return midiEvent, nil
}

// AverageLatency returns the calculated average latency.
func (lm *ListenerManager) AverageLatency() float64 {
	if lm.messageCount == 0 {
		return 0
	}
	return float64(lm.totalLatency) / float64(lm.messageCount)
}
