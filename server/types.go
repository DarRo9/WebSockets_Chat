package server

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	conn    *websocket.Conn
	server  *Server
	send    chan []byte
	address string
	name    string
}

type Server struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}
