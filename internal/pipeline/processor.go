package pipeline

import (
	"github.com/leandrodaf/pianalyze/internal/pipeline/context"
	"github.com/leandrodaf/pianalyze/internal/pipeline/stages"
	"github.com/leandrodaf/pianalyze/internal/pipeline/store"
	"go.uber.org/zap"
)

// Processor manages the execution of the pipeline by processing MIDI events through a series of stages.
type Processor struct {
	pipeline *Pipeline[context.PipelineContext, store.State]
}

// NewProcessor initializes a new pipeline processor with pre-configured stages.
// Each stage in the pipeline performs specific operations on the MIDI event context and shared state.
func NewProcessor(logger *zap.Logger) *Processor {
	state := store.NewPipelineState()
	p := NewPipeline[context.PipelineContext, store.State](state)

	// Adds stages to the pipeline in the required order
	p.AddStage(stages.NewNoteStateUpdaterStage(logger))   // Updates note state based on MIDI events
	p.AddStage(stages.NewIntervalCalculatorStage(logger)) // Calculates time intervals between events
	p.AddStage(stages.NewNoteIdentifierStage(logger))     // Identifies the current note
	p.AddStage(stages.NewChordIdentifierStage(logger))    // Identifies chords and inversions
	p.AddStage(stages.NewFinalStage(logger))              // Logs final state and sends data

	return &Processor{
		pipeline: p,
	}
}

// Process executes the pipeline stages on the provided MIDI event context.
// Returns an error if any stage in the pipeline encounters an error.
func (proc *Processor) Process(ctx *context.PipelineContext) error {
	_, err := proc.pipeline.Process(ctx)
	return err
}
