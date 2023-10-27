// message_processor.go

package controller

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"sync"
)

// MessageProcessor represents a component for processing UDP fragments and reassembling messages.
type MessageProcessor struct {
	messages map[string]*Message // Map of message IDs to incomplete messages
	mutex    sync.RWMutex
}

// Message represents a reassembled UDP message.
type Message struct {
	Data              [][]byte
	Hash              string
	IsReady           bool
	ExpectedFragments int // Number of expected fragments
	ReceivedFragments int // Number of received fragments
}

// NewMessageProcessor creates a new MessageProcessor.
func NewMessageProcessor() *MessageProcessor {
	return &MessageProcessor{
		messages: make(map[string]*Message),
	}
}

// ProcessUDPPacket handles processing of UDP packets.
func (mp *MessageProcessor) ProcessUDPPacket(packetID string, fragment []byte) {
	mp.mutex.Lock()
	defer mp.mutex.Unlock()

	message, exists := mp.messages[packetID]

	if !exists {
		// If the message doesn't exist, create a new one
		message = &Message{
			ExpectedFragments: 10, // Set the expected number of fragments
		}
		mp.messages[packetID] = message
	}

	// Append the fragment data to the message
	message.Data = append(message.Data, fragment)
	message.ReceivedFragments++

	fmt.Println("received frags : ", message.ReceivedFragments)
	fmt.Println("expected frags : ", message.ExpectedFragments)
	fmt.Println("message hash", message.Hash)

	// Check if the message is complete
	if mp.isComplete(message) {
		message.IsReady = true
		// Validate and process the complete message
		if mp.ValidateMessage(message) {
			mp.ProcessCompleteMessage(message)
		}
	}
}

// ValidateMessage validates a complete message.
func (mp *MessageProcessor) ValidateMessage(message *Message) bool {
	expectedHash := message.Hash
	actualHash := mp.computeSHA256Hash(message.Data)

	return expectedHash == actualHash
}

// ProcessCompleteMessage processes a complete message.
func (mp *MessageProcessor) ProcessCompleteMessage(message *Message) {
	// Implement the logic to process a complete message, such as saving it to a database
	fmt.Println("Processing complete message:")
	fmt.Printf("Data: %s\n", message.Data)
	fmt.Printf("SHA256 Hash: %s\n", message.Hash)
}

// isComplete checks if a message is complete based on the number of received fragments.
func (mp *MessageProcessor) isComplete(message *Message) bool {
	// Check if the message contains all expected fragments
	return message.ReceivedFragments >= message.ExpectedFragments
}

func (mp *MessageProcessor) computeSHA256Hash(data [][]byte) string {
	hash := sha256.Sum256(bytes.Join(data, []byte{}))
	return fmt.Sprintf("%x", hash)
}
