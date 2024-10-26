package listeners

import (
	"fmt"

	"github.com/leandrodaf/midi-client/pkg/pubsub"
)

// VelocityAnalyzerListener analyzes note velocity from MIDI messages and publishes the results.
type VelocityAnalyzerListener struct {
	OutputTopic string
	Hub         *pubsub.Hub
}

// ProcessMessage processes the incoming message to analyze note velocity and publishes the result.
func (l *VelocityAnalyzerListener) ProcessMessage(msg pubsub.Message) {
	velocity := msg.Data[2]

	log := fmt.Sprintf("Note velocity: %d", velocity)
	result := pubsub.Message{Data: []byte(log)}

	l.Hub.Publish(l.OutputTopic, result)
}
