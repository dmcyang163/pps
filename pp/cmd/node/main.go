package main

import (
	"fmt"
	"log"
	"pp/internal/config"
	"pp/internal/filetransfer"
	"pp/internal/message/handlers"
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

	// Create file transfer manager
	fileTransferManager := filetransfer.NewManager(cfg.DataDir)

	// Create network server
	networkServer := network.NewServer(":" + strconv.Itoa(cfg.Port))

	// Create node
	node, err := node.NewNode(cfg, networkServer) // 传入接口
	if err != nil {
		log.Fatalf("Error creating node: %v", err)
		return
	}

	// Create handlers
	pingHandler := handlers.NewPingHandler(":"+strconv.Itoa(cfg.Port), node.SendMessage)
	chatHandler := handlers.NewChatHandler()
	fileRequestHandler := handlers.NewFileRequestHandler(fileTransferManager, ":"+strconv.Itoa(cfg.Port), node.SendMessage)
	fileChunkHandler := handlers.NewFileChunkHandler(fileTransferManager)
	fileMetadataHandler := handlers.NewFileMetadataHandler(fileTransferManager)

	// Register handlers
	node.MessageRouter.RegisterHandler("ping", pingHandler)
	node.MessageRouter.RegisterHandler("chat", chatHandler)
	node.MessageRouter.RegisterHandler("file_request", fileRequestHandler)
	node.MessageRouter.RegisterHandler("file_chunk", fileChunkHandler)
	node.MessageRouter.RegisterHandler("file_metadata", fileMetadataHandler)

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
