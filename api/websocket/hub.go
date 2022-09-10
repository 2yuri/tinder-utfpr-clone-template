package websocket

import (
	"context"
	"sync"
)

type Hub struct {
	mu sync.Mutex

	clients map[string]map[*client]struct{}
	unregister chan *client
}

func (h *Hub) InsertClient(c *client, userId string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[userId]; !ok {
		h.clients[userId] = make(map[*client]struct{})
	}

	clients := h.clients[userId]
	clients[c] = struct{}{}

	h.clients[userId] = clients
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]map[*client]struct{}),
		unregister: make(chan *client),
	}
}

func (h *Hub) handleClientDelete(c *client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	c.close()
	for _, clients := range h.clients {
		for cl := range clients {
			if cl == c {
				delete(clients, c)
			}
		}
	}
}

func (h *Hub) deleteAll(){
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, clients := range h.clients {
		for cl := range clients {
			cl.close()
			delete(clients, cl)			
		}
	}
}

func (h *Hub) StartServer(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			h.deleteAll()
		case client := <- h.unregister:
			h.handleClientDelete(client)
		case event := <- Events:
			for c := range h.clients[event.UserID] {
				if c.IsAlive() {
					c.out <- event
				}
			}
		}
	}
}