// main.go

package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yuvrajsingh79/UDP-ReAssembler/pkg/controller"
)

func main() {
	// Initialize the logger with the desired log level
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)

	config := loadConfig()

	messageProcessor := controller.NewMessageProcessor()

	// Initialize the UDP ReAssembler server
	server := controller.NewServer(config, messageProcessor)

	// Start the UDP server
	server.StartUDPServer()

	// Start the HTTP server for monitoring and management
	server.StartHTTPServer()

	// Start the cleanup task for incomplete messages
	server.StartCleanupTask()

	// Block the main goroutine
	select {}
}

func loadConfig() *controller.Config {
	viper.SetConfigName("config")    // Name of the config file (without extension)
	viper.AddConfigPath("../config") // Path to the directory where config file is located
	viper.AutomaticEnv()             // Read from environment variables

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		// Handle errors when the configuration file cannot be found or read
		logrus.Fatalf("Error reading configuration: %v", err)
	}

	var config controller.Config
	err = viper.Unmarshal(&config)
	if err != nil {
		// Handle errors when the configuration cannot be unmarshaled
		logrus.Fatalf("Error unmarshaling configuration: %v", err)
	}

	return &config
}
