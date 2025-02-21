package main

import (
	"fmt"
	"log"
	"pp/internal/config"
	"pp/internal/network"
	"pp/internal/node"
	"strconv"
)

func main() {
	// Load configuration
	cfg := &config.Config{
		Port:         8080,
		MaxPeers:     10,
		SeedNodes:    []string{},
		DataDir:      "data",
		PingInterval: 10,
	}

	// Create network server
	networkServer := network.NewServer(":" + strconv.Itoa(cfg.Port))

	// Create node
	node, err := node.NewNode(cfg, networkServer) // 传入接口
	if err != nil {
		log.Fatalf("Error creating node: %v", err)
		return
	}

	// Start node
	if err := node.Start(); err != nil {
		log.Fatalf("Error starting node: %v", err)
		return
	}

	// Shutdown node on exit
	defer node.Shutdown()

	// Wait for shutdown signal (e.g., Ctrl+C)
	fmt.Println("Node running. Press Ctrl+C to exit.")
	<-make(chan struct{})
}
