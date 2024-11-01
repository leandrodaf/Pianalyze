package stages

import (
	"github.com/leandrodaf/pianalyze/internal/constants"
	"go.uber.org/zap"

	"github.com/leandrodaf/midi/sdk/contracts"
	"github.com/leandrodaf/pianalyze/internal/midi"
	"github.com/leandrodaf/pianalyze/internal/pipeline/context"
	"github.com/leandrodaf/pianalyze/internal/pipeline/store"
)

// NoteStateUpdaterStage updates the state of pressed notes based on incoming MIDI events.
type NoteStateUpdaterStage struct {
	logger *zap.Logger
}

// NewNoteStateUpdaterStage creates a new instance of NoteStateUpdaterStage with a zap logger.
func NewNoteStateUpdaterStage(logger *zap.Logger) *NoteStateUpdaterStage {
	return &NoteStateUpdaterStage{logger: logger}
}

// Process updates the pressed notes state based on the current MIDI event.
// Handles Note On and Note Off events, adjusting the state and logging key actions.
func (s *NoteStateUpdaterStage) Process(ctx *context.PipelineContext, state *store.State) error {
	event := ctx.MIDIEvent

	switch event.Command {
	case byte(contracts.NoteOn):
		if event.Velocity > 0 {
			// Adds the note to the set of pressed notes.
			state.AddNote(int(event.Note))
			s.logger.Info(constants.MsgNoteOnDetected,
				zap.String("note", midi.GetNoteName(int(event.Note))),
				zap.Int("velocity", int(event.Velocity)),
				zap.Int("command", int(event.Command)))
		} else {
			// Treats NoteOn with Velocity 0 as Note Off.
			state.RemoveNote(int(event.Note))
			s.logger.Debug(constants.MsgNoteOffViaVelocity0,
				zap.String("note", midi.GetNoteName(int(event.Note))),
				zap.Int("command", int(event.Command)))
		}
	case byte(contracts.NoteOff):
		// Removes the note from the set of pressed notes.
		state.RemoveNote(int(event.Note))
		s.logger.Info(constants.MsgNoteOffDetected,
			zap.String("note", midi.GetNoteName(int(event.Note))),
			zap.Int("command", int(event.Command)))
	default:
		// Logs unsupported MIDI commands for traffic analysis.
		s.logger.Debug(constants.MsgPipelineContextMIDI,
			zap.Int("command", int(event.Command)),
			zap.String("note", midi.GetNoteName(int(event.Note))),
			zap.Int("velocity", int(event.Velocity)))
	}

	return nil
}
