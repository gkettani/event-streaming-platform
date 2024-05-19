package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gkettani/event-streaming-platform/pkg/raft"
)

func main() {
	// Load cluster configuration
	config, err := raft.LoadClusterConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load cluster config: %v", err)
	}

	cluster := raft.NewCluster()

	for _, nodeConfig := range config.Nodes {
		cluster.AddNode(&nodeConfig)
	}

	// Wait for termination signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// Gracefully stop gRPC servers
	cluster.Stop()
}
