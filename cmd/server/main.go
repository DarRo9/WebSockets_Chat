package main

import (
	"log"

	"WebSockets_Chat/server"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	s := server.NewServer()
	logger.Info("Starting chat server", zap.Int("port", 8000))
	if err := s.Run(8000, logger); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
