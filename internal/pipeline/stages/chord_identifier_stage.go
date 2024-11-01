package stages

import (
	"github.com/leandrodaf/pianalyze/internal/constants"
	"go.uber.org/zap"

	"github.com/leandrodaf/pianalyze/internal/midi"
	"github.com/leandrodaf/pianalyze/internal/pipeline/context"
	"github.com/leandrodaf/pianalyze/internal/pipeline/store"
)

// ChordIdentifierStage identifies chords and triads based on pressed notes.
type ChordIdentifierStage struct {
	logger *zap.Logger
}

// NewChordIdentifierStage creates a new instance of ChordIdentifierStage with zap logger.
func NewChordIdentifierStage(logger *zap.Logger) *ChordIdentifierStage {
	return &ChordIdentifierStage{logger: logger}
}

// Process identifies the current chord and triad based on pressed notes and updates the pipeline context.
// Since unidentified chords can be common during live performance, they are handled without warnings.
func (s *ChordIdentifierStage) Process(ctx *context.PipelineContext, state *store.State) error {
	// Get currently pressed notes from the state.
	pressedNotes := state.GetPressedNotes()

	// Identify the chord based on pressed notes.
	chordName, inversionName, _, chordFound := midi.GetChordName(pressedNotes)
	if chordFound {
		ctx.Chord = &chordName
		ctx.Inversion = &inversionName
		s.logger.Info(constants.MsgChordAndInversionDetected, zap.String("chord", chordName), zap.String("inversion", inversionName))

		// Check if the chord is a triad.
		if midi.IsTriad(chordName) {
			ctx.Triad = &chordName
			s.logger.Info(constants.MsgTriadIdentified, zap.String("triad", chordName))
		} else {
			nonTriad := constants.NonTriad
			ctx.Triad = &nonTriad
			s.logger.Debug(constants.MsgNotTriad)
		}
	} else {
		// Set chord and triad as unknown if not identified.
		unknownChord := constants.UnknownChord
		ctx.Chord = &unknownChord
		ctx.Inversion = nil
		s.logger.Debug(constants.MsgUnknownChord)

		unknownTriad := constants.UnknownTriad
		ctx.Triad = &unknownTriad
		s.logger.Debug(constants.MsgUnknownTriad)
	}

	return nil
}
