package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/leandrodaf/midi/sdk/contracts"
	"github.com/leandrodaf/midi/sdk/midi"
	"github.com/leandrodaf/pianalyze/internal/constants"
	"github.com/leandrodaf/pianalyze/internal/pipeline"
	internalContext "github.com/leandrodaf/pianalyze/internal/pipeline/context"
	"go.uber.org/zap"
)

// Start initializes MIDI event capture and sets up a pipeline to process the captured events.
func Start() {
	logger := InitLogger()

	// Configure MIDI client with specific logging level and event filters.
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

	// Create a cancellable context for graceful shutdown handling.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channels for handling OS interrupt signals and tracking shutdown completion.
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	done := make(chan struct{})
	closeOnce := sync.Once{}

	// stopCapture handles MIDI capture shutdown, ensuring resources are released only once.
	stopCapture := func(reason string) {
		logger.Info(reason)
		if err := midiClient.Stop(); err != nil {
			logger.Error("Error stopping MIDI capture", zap.Error(err))
		}
		cancel()
		closeOnce.Do(func() { close(done) })
	}

	// Select and configure the MIDI device.
	deviceID, err := SetupDevice(ctx, midiClient)
	if err != nil {
		logger.Fatal(constants.MsgDeviceSelectionError, zap.Error(err))
		return
	}
	logger.Info(constants.MsgMIDIClientSetupSuccess, zap.Int("deviceID", deviceID))

	// Create a buffered channel for capturing MIDI events.
	eventChannel := make(chan contracts.MIDI, constants.MIDIChannelBufferSize)

	// Start capturing MIDI events.
	midiClient.StartCapture(eventChannel)

	// Initialize pipeline processor to handle MIDI events with the configured logger.
	pipelineProcessor := pipeline.NewProcessor(logger)

	var wg sync.WaitGroup

	// Goroutine for processing incoming MIDI events through the pipeline.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for event := range eventChannel {
			pipelineCtx := internalContext.NewPipelineContext(ctx, event)
			if err := pipelineProcessor.Process(pipelineCtx); err != nil {
				logger.Error(constants.MsgMIDIProcessingError, zap.Error(err))
			}
		}
	}()

	logger.Info(constants.MsgMIDIEventCaptureStarted)

	// Goroutine to handle OS interrupt signals and initiate shutdown.
	go func() {
		<-signalChan
		stopCapture("Received shutdown signal, stopping capture...")
	}()

	// Optional timeout-based shutdown mechanism (e.g., after 60 seconds).
	go func() {
		timer := time.NewTimer(60 * time.Second)
		defer timer.Stop()
		select {
		case <-timer.C:
			stopCapture("Timeout reached, stopping capture...")
		case <-done:
			// Do nothing if already shutdown.
		}
	}()

	// Wait for the shutdown signal.
	<-done

	// Close the event channel to signal end of event processing.
	close(eventChannel)

	// Wait for all events to be processed.
	wg.Wait()

	logger.Info("Shutdown complete")
}
