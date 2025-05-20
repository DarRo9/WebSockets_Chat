package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	name := "Anonymous"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:8000", Path: "/ws", RawQuery: fmt.Sprintf("name=%s", url.QueryEscape(name))}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("received: %s", message)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("read input:", err)
		}
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}

		err = c.WriteMessage(websocket.TextMessage, []byte(text))
		if err != nil {
			log.Println("write:", err)
			return
		}

		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
