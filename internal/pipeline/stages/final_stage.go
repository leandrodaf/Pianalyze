package stages

import (
	"github.com/leandrodaf/pianalyze/internal/constants"
	"go.uber.org/zap"

	"github.com/leandrodaf/pianalyze/internal/pipeline/context"
	"github.com/leandrodaf/pianalyze/internal/pipeline/store"
)

// FinalStage sends processed data to the server and logs the current state.
type FinalStage struct {
	logger *zap.Logger
}

// NewFinalStage creates a new instance of FinalStage with zap logger.
func NewFinalStage(logger *zap.Logger) *FinalStage {
	return &FinalStage{logger: logger}
}

// Process sends processed data to the server and logs the pipeline context and shared state.
// In a real implementation, this function would contain server communication logic.
func (s *FinalStage) Process(ctx *context.PipelineContext, state *store.State) error {
	// Helper function to safely dereference string pointers.
	getString := func(s *string) string {
		if s == nil {
			return "nil"
		}
		return *s
	}

	// Logs core details of the current MIDI event from the pipeline context.
	s.logger.Info(constants.MsgPipelineContextMIDI,
		zap.Int("command", int(ctx.MIDIEvent.Command)),
		zap.Int("note", int(ctx.MIDIEvent.Note)),
		zap.Int("velocity", int(ctx.MIDIEvent.Velocity)),
		zap.Uint64("timestamp", ctx.MIDIEvent.Timestamp))

	// Logs additional details in the pipeline context for debugging purposes.
	s.logger.Debug(constants.MsgPipelineAdditionalDetails,
		zap.Uint64("interval", ctx.Interval),
		zap.String("currentKey", getString(ctx.CurrentKey)),
		zap.String("triad", getString(ctx.Triad)),
		zap.String("chord", getString(ctx.Chord)),
		zap.String("inversion", getString(ctx.Inversion)))

	// Logs the current state with pressed notes and last note time.
	s.logger.Info(constants.MsgStatePressedNotes, zap.Any("pressedNotes", state.GetPressedNotes()))
	s.logger.Debug(constants.MsgStateLastNoteTime, zap.Uint64("lastNoteTime", state.GetLastNoteTime()))

	// Placeholder for server communication logic:
	// err := sendToServer(ctx, state)
	// if err != nil {
	//     s.logger.Error("Failed to send data to server", zap.Error(err))
	//     return err
	// }

	return nil
}
