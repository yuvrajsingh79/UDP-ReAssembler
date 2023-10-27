// storage.go

package controller

import (
	"sync"
)

// Storage is responsible for managing storage of reassembled messages.
type Storage struct {
	messages map[string]*Message
	mu       sync.RWMutex
}

// NewStorage creates a new Storage instance.
func NewStorage() *Storage {
	return &Storage{
		messages: make(map[string]*Message),
	}
}

// AddMessage adds a reassembled message to the storage.
func (s *Storage) AddMessage(packetID string, message *Message) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.messages[packetID] = message
}

// GetMessage retrieves a reassembled message from the storage.
func (s *Storage) GetMessage(packetID string) *Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.messages[packetID]
}

// RemoveMessage removes a reassembled message from the storage.
func (s *Storage) RemoveMessage(packetID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.messages, packetID)
}
