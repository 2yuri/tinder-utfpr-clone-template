package websocket

import (
	"log"
	"time"

	ws "github.com/gorilla/websocket"
)

type client struct {
	alive bool
	hub *Hub
	conn *ws.Conn

	out chan interface{}
}

func NewClient(conn *ws.Conn, hub *Hub) *client {
	c := &client{
		conn: conn, 
		hub: hub, 
		out: make(chan interface{}, 10),
		alive: true,
	}

	go c.loopIn()
	go c.loopOut()

	return c
}

func (c *client) close() {
	if c.alive {
		c.conn.Close()
		c.alive = false
		close(c.out)
		for len(c.out) > 0 {
			<-c.out
		}
	}
}

func (c *client) loopOut() {
	for {
		select {
		case m := <-c.out:
			err := c.conn.WriteJSON(m)
			if err != nil {
				c.close()
				c.hub.unregister <- c
			}
		}
	}
}

func (c *client) loopIn() {
	defer func() {
		c .close()
		c.hub.unregister <- c
	}()

	for {
		messageType, _, err := c.conn.ReadMessage()
		if err != nil {
			if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway, ws.CloseAbnormalClosure){
				log.Println(err.Error())
			}

			break
		}

		switch messageType {
		case ws.CloseMessage:
			return
		case ws.PingMessage:
			c.conn.WriteControl(ws.PongMessage, nil, time.Now().Add(time.Second * 5))
		}
	}
}

func (c *client) IsAlive() bool {
	return c.alive
}