package server

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	activeConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "websocket_active_connections",
		Help: "The total number of active WebSocket connections",
	})

	messagesReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "websocket_messages_received_total",
		Help: "The total number of messages received",
	})

	messagesSent = promauto.NewCounter(prometheus.CounterOpts{
		Name: "websocket_messages_sent_total",
		Help: "The total number of messages sent",
	})
)

func (s *Server) updateMetrics() {
	activeConnections.Set(float64(len(s.clients)))
}

func incrementMessagesReceived() {
	messagesReceived.Inc()
}

func incrementMessagesSent() {
	messagesSent.Inc()
}
