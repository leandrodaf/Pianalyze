package listeners

import (
	"encoding/binary"
	"time"

	"github.com/leandrodaf/midi-client/internal/contracts/logger"
	"github.com/leandrodaf/midi-client/pkg/pubsub"
)

// ListenerManager manages a set of listeners and processes events for them.
type ListenerManager struct {
	hub          *pubsub.Hub
	topicName    string
	listeners    []Listener
	logger       logger.Logger
	totalLatency int64
	messageCount int64
}

// NewListenerManager creates a new instance of ListenerManager.
func NewListenerManager(hub *pubsub.Hub, topicName string, logger logger.Logger) *ListenerManager {
	return &ListenerManager{
		hub:       hub,
		topicName: topicName,
		logger:    logger,
	}
}

// Register adds a listener to the manager.
func (lm *ListenerManager) Register(listener Listener) {
	lm.listeners = append(lm.listeners, listener)
}

// Start begins listening and processing events for all registered listeners.
func (lm *ListenerManager) Start() {
	eventChannel := lm.hub.Subscribe(lm.topicName)
	go func() {
		defer lm.hub.Unsubscribe(lm.topicName, eventChannel)
		for msg := range eventChannel {
			lm.processMessage(msg)
		}
	}()
}

// processMessage distributes the MIDI event to all registered listeners.
func (lm *ListenerManager) processMessage(msg pubsub.Message) {
	for _, listener := range lm.listeners {
		measureEventLatency(msg, lm)
		listener.ProcessMessage(msg)
	}
}

// measureEventLatency calculates and logs the latency of an event.
func measureEventLatency(msg pubsub.Message, lm *ListenerManager) {
	eventTimestamp := int64(binary.BigEndian.Uint64(msg.Data[:8]))
	receiveTimestamp := time.Now().UnixNano()
	latency := receiveTimestamp - eventTimestamp

	lm.totalLatency += latency
	lm.messageCount++

	lm.logger.Info("Listener latency", lm.logger.Field().Int64("latency_ns", latency))
}

// AverageLatency returns the calculated average latency.
func (lm *ListenerManager) AverageLatency() float64 {
	if lm.messageCount == 0 {
		return 0
	}
	return float64(lm.totalLatency) / float64(lm.messageCount)
}
