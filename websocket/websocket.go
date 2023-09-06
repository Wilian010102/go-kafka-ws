package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
)

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer closeConnection(conn) // Close the WebSocket connection when done

	addClient(conn) // Add the client to the clients map

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
	}
}

// StartWebSocketServer starts the WebSocket server
func StartWebSocketServer() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading connection:", err)
			return
		}
		defer closeConnection(conn) // Close the WebSocket connection when done

		addClient(conn) // Add the client to the clients map

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}
		}
	})

	http.ListenAndServe(":9999", nil)
}

// BroadcastMessage sends a message to all connected clients
func BroadcastMessage(message string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Error sending message to client:", err)
		}
	}
}

// SendWebSocketUpdate sends a WebSocket update
func SendWebSocketUpdate(message string) {
	go BroadcastMessage(message)
}

func closeConnection(conn *websocket.Conn) {
	conn.Close()

	clientsMu.Lock()
	delete(clients, conn) // Remove the client from the clients map
	clientsMu.Unlock()
}

func addClient(conn *websocket.Conn) {
	clientsMu.Lock()
	clients[conn] = true // Add the client to the clients map
	clientsMu.Unlock()
}
