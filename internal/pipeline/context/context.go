package context

import (
	"context"

	"github.com/leandrodaf/midi/sdk/contracts"
	"github.com/leandrodaf/pianalyze/internal/constants"
)

// PipelineContext is a custom context that embeds context.Context
// and includes fields for MIDI event data and music-related information.
type PipelineContext struct {
	context.Context
	MIDIEvent  contracts.MIDI // Current MIDI event data
	Interval   uint64         // Time interval between consecutive events
	CurrentKey *string        // Detected current note or key
	Triad      *string        // Identified triad, if applicable
	Chord      *string        // Identified full chord, if applicable
	Inversion  *string        // Chord inversion, if applicable
}

// NewPipelineContext initializes a new PipelineContext with a parent context and a MIDI event.
// Default values are assigned to musical fields, providing initial references for key, triad, chord, and inversion.
func NewPipelineContext(ctx context.Context, event contracts.MIDI) *PipelineContext {
	defaultKey := constants.DefaultKey
	defaultTriad := constants.DefaultTriad
	defaultChord := constants.DefaultChord
	defaultInversion := constants.DefaultInversion
	return &PipelineContext{
		Context:    ctx,
		MIDIEvent:  event,
		Interval:   0,
		CurrentKey: &defaultKey,
		Triad:      &defaultTriad,
		Chord:      &defaultChord,
		Inversion:  &defaultInversion,
	}
}
