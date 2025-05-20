package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) Run(port int, logger *zap.Logger) error {
	http.HandleFunc("/ws", s.handleWebSocket)

	go s.run(logger)

	addr := fmt.Sprintf(":%d", port)
	logger.Info("Starting server", zap.String("address", addr))
	return http.ListenAndServe(addr, nil)
}

func (s *Server) run(logger *zap.Logger) {
	for {
		select {
		case client := <-s.register:
			s.clients[client] = true
			welcome := client.name
			if welcome == "" {
				welcome = client.address
			}
			message := fmt.Sprintf("You are %s", welcome)
			client.send <- []byte(message)
			logger.Info("New client connected", zap.String("address", client.address), zap.String("name", client.name))

		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
				left := client.name
				if left == "" {
					left = client.address
				}
				message := fmt.Sprintf("%s has left", left)
				s.broadcast <- []byte(message)
				logger.Info("Client disconnected", zap.String("address", client.address), zap.String("name", client.name))
			}

		case message := <-s.broadcast:
			for client := range s.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(s.clients, client)
				}
			}
		}
	}
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	name := r.URL.Query().Get("name")
	client := &Client{
		conn:    conn,
		server:  s,
		send:    make(chan []byte, 256),
		address: r.RemoteAddr,
		name:    name,
	}

	s.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.server.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		sender := c.name
		if sender == "" {
			sender = c.address
		}
		formattedMessage := fmt.Sprintf("%s: %s", sender, string(message))
		c.server.broadcast <- []byte(formattedMessage)
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}
