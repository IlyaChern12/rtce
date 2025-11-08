package realtime

import (
	"encoding/json"
	"log"

	"github.com/IlyaChern12/rtce/internal/models"
	"github.com/gorilla/websocket"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
}

func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("invalid message:", err)
			continue
		}

		docID, ok := msg["documentId"].(string)
		if !ok {
			log.Println("missing documentId")
			continue
		}

		doc, ok := hub.documents[docID]
		if !ok {
			// создаем новый документ в Hub, если нет
			doc = models.NewDocumentCRDT(docID)
			hub.documents[docID] = doc
		}

		switch msg["action"] {
		case "insert":
			charMap := msg["char"].(map[string]interface{})
			char := models.Char{
				ID:     charMap["id"].(string),
				Value:  charMap["value"].(string),
				PrevID: charMap["prevId"].(string),
			}
			doc.Insert(char)
		case "delete":
			charID := msg["charId"].(string)
			doc.Delete(charID)
		default:
			log.Println("unknown action:", msg["action"])
			continue
		}

		// рассылаем всем остальным клиентам
		hub.broadcast <- BroadcastMessage{
			Sender: c,
			Data:   message,
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Printf("writePump error: %v", err)
			break
		}
	}
}