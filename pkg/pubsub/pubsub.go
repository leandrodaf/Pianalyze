package pubsub

import (
	"sync"
)

// Defining the message type for binary data or lightweight structs
type Message struct {
	Data []byte
}

// Hub is the main structure of the pub-sub system
type Hub struct {
	mu     sync.RWMutex
	topics map[string]*Topic
}

// Topic represents a topic with its subscribers
type Topic struct {
	subscribers sync.Map // key: chan Message, value: struct{}
}

// NewHub creates a new instance of Hub
func NewHub() *Hub {
	return &Hub{
		topics: make(map[string]*Topic),
	}
}

// Subscribe allows a client to subscribe to a topic
func (h *Hub) Subscribe(topicName string) chan Message {
	ch := make(chan Message, 100) // Buffered channel to avoid blocking
	topic := h.getOrCreateTopic(topicName)
	topic.addSubscriber(ch)
	return ch
}

// Publish sends a message to a specific topic
func (h *Hub) Publish(topicName string, msg Message) {
	topic := h.getOrCreateTopic(topicName)
	topic.publish(msg)
}

// Unsubscribe removes a subscriber from a topic
func (h *Hub) Unsubscribe(topicName string, ch chan Message) {
	h.mu.RLock()
	topic, exists := h.topics[topicName]
	h.mu.RUnlock()
	if exists {
		topic.removeSubscriber(ch)
		if topic.subscriberCount() == 0 {
			h.mu.Lock()
			delete(h.topics, topicName)
			h.mu.Unlock()
		}
	}
}

// getOrCreateTopic retrieves or creates a topic
func (h *Hub) getOrCreateTopic(topicName string) *Topic {
	h.mu.RLock()
	topic, exists := h.topics[topicName]
	h.mu.RUnlock()
	if exists {
		return topic
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	// Check again to avoid race condition
	if topic, exists = h.topics[topicName]; exists {
		return topic
	}

	topic = &Topic{}
	h.topics[topicName] = topic
	return topic
}

// Methods of the Topic struct

func (t *Topic) addSubscriber(ch chan Message) {
	t.subscribers.Store(ch, struct{}{})
}

func (t *Topic) removeSubscriber(ch chan Message) {
	t.subscribers.Delete(ch)
	close(ch)
}

func (t *Topic) subscriberCount() int {
	count := 0
	t.subscribers.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

func (t *Topic) publish(msg Message) {
	t.subscribers.Range(func(key, value interface{}) bool {
		ch := key.(chan Message)
		select {
		case ch <- msg:
		default:
			// Optional: handle slow subscribers
		}
		return true
	})
}
