package stages

import (
	"go.uber.org/zap"

	"github.com/leandrodaf/pianalyze/internal/constants"
	"github.com/leandrodaf/pianalyze/internal/midi"
	"github.com/leandrodaf/pianalyze/internal/pipeline/context"
	"github.com/leandrodaf/pianalyze/internal/pipeline/store"
)

// NoteIdentifierStage identifies the current note based on the pressed notes and updates the pipeline context.
type NoteIdentifierStage struct {
	logger *zap.Logger
}

// NewNoteIdentifierStage creates a new instance of NoteIdentifierStage with a zap logger.
func NewNoteIdentifierStage(logger *zap.Logger) *NoteIdentifierStage {
	return &NoteIdentifierStage{logger: logger}
}

// Process identifies the current note based on the last pressed note and updates the context with it.
// If no notes are currently pressed, it clears the current key in the context.
func (s *NoteIdentifierStage) Process(ctx *context.PipelineContext, state *store.State) error {
	// Retrieves currently pressed notes from the state.
	pressedNotes := state.GetPressedNotes()

	// Checks if there are pressed notes to identify the current note.
	if len(pressedNotes) > 0 {
		// Identifies the last pressed note.
		lastNote := pressedNotes[len(pressedNotes)-1]
		noteName := midi.GetNoteName(lastNote)
		ctx.CurrentKey = &noteName
		s.logger.Info(constants.MsgStatePressedNotes, zap.String("note", noteName))

	} else {
		// If no notes are pressed, clears the CurrentKey in the context.
		ctx.CurrentKey = nil
		s.logger.Debug(constants.MsgNoPreviousEvent)
	}

	return nil
}
