package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func TestServer(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	server := NewServer()
	go server.Run(8000, logger)

	time.Sleep(time.Second)

	ts := httptest.NewServer(http.HandlerFunc(server.handleWebSocket))
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", wsURL, err)
	}
	defer ws.Close()

	_, message, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("could not read message: %v", err)
	}
	if !strings.Contains(string(message), "You are") {
		t.Errorf("expected welcome message, got %s", message)
	}

	err = ws.WriteMessage(websocket.TextMessage, []byte("test message"))
	if err != nil {
		t.Fatalf("could not send message: %v", err)
	}

	_, message, err = ws.ReadMessage()
	if err != nil {
		t.Fatalf("could not read message: %v", err)
	}
	if !strings.Contains(string(message), "test message") {
		t.Errorf("expected message to contain 'test message', got %s", message)
	}
}
