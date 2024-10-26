package channel

import (
	"github.com/leandrodaf/pianalyze/internal/entity"
)

// MIDIEventPublisher defines an interface for publishing and managing MIDI events.
type MIDIEventPublisher interface {
	Publish(event entity.MIDI)  // Publishes a MIDI event to the channel.
	Events() <-chan entity.MIDI // Returns the MIDI event channel for reading.
	Shutdown()                  // Shuts down the event channel.
	GetChannel() chan entity.MIDI
}
