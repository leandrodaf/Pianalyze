package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/leandrodaf/midi/sdk/contracts"
	"github.com/leandrodaf/midi/sdk/midi"
	"github.com/leandrodaf/pianalyze/internal/constants"
	"github.com/leandrodaf/pianalyze/internal/pipeline"
	internalContext "github.com/leandrodaf/pianalyze/internal/pipeline/context"
	"go.uber.org/zap"
)

// Start initializes the MIDI event capture process and sets up the pipeline processor.
func Start() {
	logger := InitLogger()

	// Sets up the MIDI client with specific logging level and event filter.
	midiClient, err := midi.NewMIDIClient(
		contracts.WithLogLevel(contracts.InfoLevel),
		contracts.WithMIDIEventFilter(contracts.MIDIEventFilter{
			Commands: []contracts.MIDICommand{contracts.NoteOn, contracts.NoteOff},
		}),
	)
	if err != nil {
		logger.Error(constants.MsgMIDIClientSetupError, zap.Error(err))
		return
	}

	// Selects and configures the MIDI device.
	deviceID, err := SetupDevice(midiClient)
	if err != nil {
		logger.Fatal(constants.MsgDeviceSelectionError, zap.Error(err))
		return
	}
	logger.Info(constants.MsgMIDIClientSetupSuccess, zap.Int("deviceID", deviceID))

	// Sets up a buffered channel for capturing MIDI events.
	eventChannel := make(chan contracts.MIDI, constants.MIDIChannelBufferSize)
	midiClient.StartCapture(eventChannel)

	// Initializes the pipeline processor with the configured logger.
	pipelineProcessor := pipeline.NewProcessor(logger)

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	// Channel to listen for OS signals (like CTRL+C)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM) // Capture CTRL+C and SIGTERM

	// Goroutine to process incoming MIDI events.
	wg.Add(1)
	go func() {
		defer wg.Done() // Decrement the WaitGroup counter when the goroutine completes
		for event := range eventChannel {
			ctx := internalContext.NewPipelineContext(ctx, event)
			if err := pipelineProcessor.Process(ctx); err != nil {
				logger.Error(constants.MsgMIDIProcessingError, zap.Error(err))
			}
		}
	}()

	logger.Info(constants.MsgMIDIEventCaptureStarted)

	// Wait for a signal
	<-signalChan
	cancel() // Cancel the context to signal goroutines to stop

	logger.Info("Received shutdown signal, stopping...")

	// Close the event channel to stop the processing goroutine
	close(eventChannel) // Close the event channel to stop the processing goroutine
	wg.Wait()           // Wait for the processing goroutine to finish

	defer func() { _ = midiClient.Stop() }()

	logger.Info("Shutdown complete")
}
