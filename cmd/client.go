package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

var stopCh chan error

const TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjI5MDAzNDYsImlhdCI6MTY2MjgxMzk0NiwiaXNzIjoidGluZGVyLWlzc3VlciIsInN1YiI6ImRiZmIzZmQyLTMxMDUtMTFlZC05NTgwLTAyNDJhYzFkMDAwMiJ9.FB0WNcZ4vEzBnCYPQTUfWoePk-1eY95KGhSzz8r0TT8"

func main() {
	flag.Parse()
	log.SetFlags(0)
	stopCh = make(chan error)

	u := url.URL{Scheme: "ws", Host: "localhost:51000", Path: "/api/v1/subscribe"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{
		"Authorization": []string{"Bearer " + TOKEN},
	})
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	consume(c)
}

func consume(c *websocket.Conn) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			stopCh <- err
			return
		}
		log.Printf("%v - recv: %s", time.Now(), message)
	}
}
