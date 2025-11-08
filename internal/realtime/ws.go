package realtime

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// принимает соединение и делает эхо
func WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade websocket:", err)
		return
	}
	defer conn.Close()

	log.Println("New WS connection established")

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		log.Printf("Received: %s", message)

		// эхо обратно клиенту
		if err := conn.WriteMessage(mt, message); err != nil {
			log.Println("write error:", err)
			break
		}
	}
}
