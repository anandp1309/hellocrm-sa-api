package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for this example, restrict in production
	},
}

type Hub struct {
	clients    map[*websocket.Conn]bool
	clientsMux sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (h *Hub) HandleWebSocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return err
	}
	defer ws.Close()

	h.clientsMux.Lock()
	h.clients[ws] = true
	h.clientsMux.Unlock()

	defer func() {
		h.clientsMux.Lock()
		delete(h.clients, ws)
		h.clientsMux.Unlock()
	}()

	for {
		// Read message
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("WebSocket error: %v", err)
			break
		}
		
		log.Printf("Received WS message: %s", msg)

		// Echo back for demonstration
		err = ws.WriteMessage(websocket.TextMessage, []byte("Echo: "+string(msg)))
		if err != nil {
			log.Printf("WebSocket write error: %v", err)
			break
		}
	}
	return nil
}
