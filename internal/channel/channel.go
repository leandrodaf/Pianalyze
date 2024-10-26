package channel

import (
	"sync/atomic"

	"github.com/leandrodaf/pianalyze/internal/contracts/channel"
	"github.com/leandrodaf/pianalyze/internal/entity"
)

// EventChannel manages the publication and retrieval of MIDI events.
type EventChannel struct {
	eventChan chan entity.MIDI
	closed    int32
}

// NewEventChannel creates a new EventChannel with a specified buffer size.
func NewEventChannel(bufferSize int) channel.MIDIEventPublisher {
	return &EventChannel{
		eventChan: make(chan entity.MIDI, bufferSize),
	}
}

// Publish safely publishes a MIDI event to the channel in a non-blocking manner.
func (ec *EventChannel) Publish(event entity.MIDI) {
	if atomic.LoadInt32(&ec.closed) == 1 {
		return // Channel is closed, discard the event or handle as needed
	}
	select {
	case ec.eventChan <- event:
	default:
		// Buffer is full, discard the event or handle as needed
	}
}

// Events returns the MIDI event channel for reading.
func (ec *EventChannel) Events() <-chan entity.MIDI {
	return ec.eventChan
}

// Shutdown safely closes the event channel.
func (ec *EventChannel) Shutdown() {
	if atomic.CompareAndSwapInt32(&ec.closed, 0, 1) {
		close(ec.eventChan)
	}
}

// GetChannel returns the channel for MIDI event reading.
func (ec *EventChannel) GetChannel() chan entity.MIDI {
	return ec.eventChan
}
