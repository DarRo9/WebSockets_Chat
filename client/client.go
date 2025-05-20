package main

import (
	"bufio"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	name := ""
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	u := url.URL{Scheme: "ws", Host: "localhost:8000", Path: "/ws"}
	q := u.Query()
	if name != "" {
		q.Set("name", name)
	}
	u.RawQuery = q.Encode()

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

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		err := c.WriteMessage(websocket.TextMessage, scanner.Bytes())
		if err != nil {
			log.Println("write:", err)
			return
		}
	}

	select {
	case <-done:
		return
	case <-interrupt:
		log.Println("interrupt")

		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("write close:", err)
			return
		}
		select {
		case <-done:
		}
	}
}
