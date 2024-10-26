package cmd

import (
	"fmt"
	"time"
)

// Start initializes the MIDI event capture process.
func Start() {
	logger, err := SetupLogger()
	if err != nil {
		fmt.Println("Error setting up logger:", err)
		return
	}

	midiClient, err := SetupMIDIClient(logger)
	if err != nil {
		logger.Error("Error setting up MIDI client", logger.Field().String("error", err.Error()))
		return
	}

	deviceID, err := SetupDevice(midiClient)
	if err != nil {
		logger.Fatal("Error selecting MIDI device", logger.Field().String("error", err.Error()))
		return
	}
	logger.Info("MIDI device successfully selected", logger.Field().Int("deviceID", deviceID))

	eventChannel := SetupEventChannel(100)
	hub, topicName, processedTopicName := SetupHub()

	defer eventChannel.Shutdown()
	defer midiClient.Stop()

	midiClient.StartCapture(eventChannel.GetChannel())
	StartEventPublishing(eventChannel, hub, topicName)

	latencyInterval := 5 * time.Second

	listenerManager := SetupListenerManager(hub, topicName, processedTopicName, logger, latencyInterval)
	listenerManager.Start()

	fmt.Println("Capturing MIDI events. Press Ctrl+C to stop.")
	select {}
}
