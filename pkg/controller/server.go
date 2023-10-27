package controller

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

// UDPConn is a package-level variable to hold the UDP connection.
var UDPConn *net.UDPConn

// Config represents the configuration for the UDP ReAssembler server.
type Config struct {
	UDPPort                          int
	UDPPacketBufferSize              int
	HTTPListenAddress                string
	IncompleteMessageCleanupInterval time.Duration
	LogLevel                         string
	FragmentCacheTimeout             time.Duration
}

// Server represents the UDP ReAssembler server.
type Server struct {
	config           *Config
	messageProcessor *MessageProcessor
}

// NewServer creates a new Server instance.
func NewServer(config *Config, messageProcessor *MessageProcessor) *Server {
	return &Server{
		config:           config,
		messageProcessor: messageProcessor,
	}
}

// // StartUDPServer starts the UDP server.
// func (s *Server) StartUDPServer() {
// 	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: s.config.UDPPort})
// 	if err != nil {
// 		logger.Fatalf("Error starting UDP server: %v", err)
// 	}
// 	defer conn.Close()

// 	logger.Infof("UDP server listening on port %d", s.config.UDPPort)

// 	// Handle incoming UDP packets
// 	go func() {
// 		for {
// 			buffer := make([]byte, s.config.UDPPacketBufferSize)
// 			n, addr, err := conn.ReadFromUDP(buffer)
// 			if err != nil {
// 				logger.Errorf("Error reading UDP packet: %v", err)
// 				continue
// 			}

// 			packetData := buffer[:n]

// 			// Process the received packet
// 			s.messageProcessor.ProcessUDPPacket(addr.String(), packetData)
// 		}
// 	}()
// }

// StartUDPServer starts the UDP server.
func (s *Server) StartUDPServer() {
	var err error
	UDPConn, err = net.ListenUDP("udp", &net.UDPAddr{Port: s.config.UDPPort})
	if err != nil {
		logger.Fatalf("Error starting UDP server: %v", err)
	}
	fmt.Println("UDP message : ", s.messageProcessor.messages)
	logger.Infof("UDP server listening on  port %d ", s.config.UDPPort)

	// Handle incoming UDP packets
	go func() {
		for {
			buffer := make([]byte, s.config.UDPPacketBufferSize)
			n, addr, err := UDPConn.ReadFromUDP(buffer)
			if err != nil {
				logger.Errorf("Error reading UDP packet: %v", err)
				continue
			}

			packetData := buffer[:n]

			// Process the received packet
			s.messageProcessor.ProcessUDPPacket(addr.String(), packetData)
		}
	}()
}

// StartHTTPServer starts the HTTP server for monitoring and management.
func (s *Server) StartHTTPServer() {
	router := mux.NewRouter()
	router.HandleFunc("/status", s.statusHandler)

	httpServer := &http.Server{
		Addr:    s.config.HTTPListenAddress,
		Handler: router,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Fatalf("HTTP server error: %v", err)
		}
	}()
}

func (s *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "UDP ReAssembler is running.")
}

// StartCleanupTask starts the periodic cleanup of incomplete messages.
func (s *Server) StartCleanupTask() {
	go func() {
		ticker := time.NewTicker(s.config.IncompleteMessageCleanupInterval)
		for range ticker.C {
			s.cleanupIncompleteMessages()
		}
	}()
}

func (s *Server) cleanupIncompleteMessages() {
	// Implement the cleanup logic to remove incomplete messages based on the project requirements.
	// In this example, we remove messages that haven't received all expected fragments within a certain time period.
	// You should adapt this logic to your specific project requirements.

	for packetID, message := range s.messageProcessor.messages {
		if !message.IsReady {
			// You may need to add your own logic here to determine when a message should be considered incomplete.
			// For instance, you can check the number of received fragments and the expected number of fragments.
			if message.ReceivedFragments < message.ExpectedFragments {
				logger.Infof("Cleaning up incomplete message with ID: %s", packetID)
				delete(s.messageProcessor.messages, packetID)
			}
		}
	}
}
