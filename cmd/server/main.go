package main

import (
	"flag"
	"log"

	"websocket-chat/server"

	"go.uber.org/zap"
)

func main() {
	port := flag.Int("port", 8000, "port to listen on")
	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	s := server.NewServer()
	logger.Info("Starting chat server", zap.Int("port", *port))

	if err := s.Run(*port, logger); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
