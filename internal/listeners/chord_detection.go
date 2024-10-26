package listeners

import (
	"fmt"

	"github.com/leandrodaf/pianalyze/pkg/pubsub"
)

// ChordDetectionListener detects chords based on MIDI messages and publishes the results.
type ChordDetectionListener struct {
	OutputTopic string
	Hub         *pubsub.Hub
}

// ProcessMessage processes the incoming message to detect chords and publishes the result.
func (l *ChordDetectionListener) ProcessMessage(msg pubsub.Message) {
	note := msg.Data[1]
	velocity := msg.Data[2]

	chordName := fmt.Sprintf("Chord detected for note %d with velocity %d", note, velocity)
	result := pubsub.Message{Data: []byte(chordName)}

	l.Hub.Publish(l.OutputTopic, result)
}
