package realtime

import "github.com/IlyaChern12/rtce/internal/models"

type BroadcastMessage struct {
	Sender *Client
	Data   []byte
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan BroadcastMessage
	register   chan *Client
	unregister chan *Client

	documents map[string]*models.Document
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        broadcast:  make(chan BroadcastMessage),
        documents:  make(map[string]*models.Document),
    }
}

func (h *Hub) Run() {
	for {
	select {
	case client := <-h.register:
		h.clients[client] = true
	case client := <-h.unregister:
		if _, ok := h.clients[client]; ok {
			delete(h.clients, client)
			close(client.send)
		}
	case msg := <-h.broadcast:
		for client := range h.clients {
			if client != msg.Sender {
				select {
				case client.send <- msg.Data:
				default:
					delete(h.clients, client)
					close(client.send)
				}
			}
		}
	}
}
}
