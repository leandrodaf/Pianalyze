package stages

import (
	"github.com/leandrodaf/pianalyze/internal/constants"
	"go.uber.org/zap"

	"github.com/leandrodaf/pianalyze/internal/pipeline/context"
	"github.com/leandrodaf/pianalyze/internal/pipeline/store"
)

// IntervalCalculatorStage calculates the time interval between consecutive MIDI events.
type IntervalCalculatorStage struct {
	logger *zap.Logger
}

// NewIntervalCalculatorStage creates a new instance of IntervalCalculatorStage with zap logger.
func NewIntervalCalculatorStage(logger *zap.Logger) *IntervalCalculatorStage {
	return &IntervalCalculatorStage{logger: logger}
}

// Process calculates the time interval between the current and previous MIDI events and updates the context with this interval.
// If there is no previous event, it sets the interval to zero.
func (s *IntervalCalculatorStage) Process(ctx *context.PipelineContext, state *store.State) error {
	// Retrieves the timestamp of the current MIDI event.
	currentTime := ctx.MIDIEvent.Timestamp

	// Retrieves the timestamp of the last MIDI event.
	lastTime := state.GetLastNoteTime()

	if lastTime > 0 {
		// Calculates the time interval between the current and last event.
		ctx.Interval = currentTime - lastTime
		s.logger.Info(constants.MsgIntervalCalculated,
			zap.Uint64("interval", ctx.Interval))
	} else {
		// If no previous event exists, sets interval to zero.
		ctx.Interval = 0
		s.logger.Debug(constants.MsgNoPreviousEvent)
	}

	// Updates the state with the current timestamp as the last note time.
	state.UpdateLastNoteTime(currentTime)

	return nil
}
