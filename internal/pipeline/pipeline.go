package pipeline

// Stage represents a stage in the pipeline that processes `TContext` using `TState`.
// Each stage performs a specific operation on the context with access to shared state.
type Stage[TContext any, TState any] interface {
	Process(ctx *TContext, state *TState) error
}

// Pipeline represents a sequence of stages that process data of type `TContext` with shared `TState`.
// The pipeline manages the execution of each stage in a specific order, passing along context and state.
type Pipeline[TContext any, TState any] struct {
	stages []Stage[TContext, TState]
	state  *TState
}

// NewPipeline creates a new pipeline with the given shared state.
// The pipeline starts with an empty sequence of stages, which can be added as needed.
func NewPipeline[TContext any, TState any](state *TState) *Pipeline[TContext, TState] {
	return &Pipeline[TContext, TState]{
		stages: []Stage[TContext, TState]{},
		state:  state,
	}
}

// AddStage adds a stage to the pipeline.
// Stages are executed in the order they are added.
func (p *Pipeline[TContext, TState]) AddStage(stage Stage[TContext, TState]) {
	p.stages = append(p.stages, stage)
}

// Process executes the pipeline by processing the given `TContext` through each stage in sequence.
// If any stage returns an error, the processing stops, and the error is returned.
// Returns the processed context or nil if the input context is nil.
func (p *Pipeline[TContext, TState]) Process(ctx *TContext) (*TContext, error) {
	if ctx == nil {
		return nil, nil
	}

	for _, stage := range p.stages {
		if err := stage.Process(ctx, p.state); err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}
