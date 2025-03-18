package handlers

import (
	"db"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"sync"

	"github.com/gorilla/websocket"
)

// Storing the clients connected to the server in a map
var clients = make(map[string]*websocket.Conn)

// Mutex to be able to lock before writing to the map
var mu sync.Mutex

// Upgrader to upgrade the HTTP connection to a WebSocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow WebSocket connections from http://localhost:8080
		allowedOrigins := []string{
			"https://localhost:8080",
		}
		// Get the website link (ex: http://localhost:8080)
		origin := r.Header.Get("Origin")
		// Return true if the link is contained in the allowed links
		return slices.Contains(allowedOrigins, origin)
	},
}

// Struct that will store the content of the json sent by the javascript
type Message struct {
	Type     string `json:"type"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}

// Handler that will upgrade the HTTP connection to a WebSocket connection and listen for messages
func HandleConnection(w http.ResponseWriter, r *http.Request) {
	// Retrieving the cookie for the uuid
	cookie, errCookie := r.Cookie("session_id")
	if errCookie != nil {
		fmt.Println("Error getting cookie:", errCookie)
		return
	}
	// Getting the username with the UUID stored in the cookie
	username := db.UserNicknameWithUUID(cookie.Value)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()

	mu.Lock()
	clients[username] = conn
	mu.Unlock()

	fmt.Println(username, "connected")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(username, "disconnected")
			mu.Lock()
			delete(clients, username)
			mu.Unlock()
			break
		}

		var receivedMsg Message
		err = json.Unmarshal(msg, &receivedMsg)
		if err != nil {
			fmt.Println("Invalid JSON:", err)
			continue
		}

		if receivedMsg.Type == "private_message" {
			// sendPrivateMessage(receivedMsg)
		}
	}
}

// package handlers

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"sync"
// 	"time"

// 	"github.com/gorilla/websocket"
// )

// // Client represents a connected WebSocket client
// type Client struct {
// 	ID       string
// 	Conn     *websocket.Conn
// 	Send     chan []byte
// 	Username string
// }

// // ClientManager manages the connected clients
// type ClientManager struct {
// 	Clients    map[string]*Client
// 	Register   chan *Client
// 	Unregister chan *Client
// 	Broadcast  chan []byte
// 	PrivateMsg chan PrivateMessage
// 	mutex      sync.Mutex
// }

// // PrivateMessage contains a private message and its recipient
// type PrivateMessage struct {
// 	To      string `json:"to"`
// 	From    string `json:"from"`
// 	Content string `json:"content"`
// }

// // Message is the structure for all messages
// type Message struct {
// 	Type    string `json:"type"`
// 	From    string `json:"from"`
// 	To      string `json:"to,omitempty"`
// 	Content string `json:"content"`
// }

// var manager = ClientManager{
// 	Clients:    make(map[string]*Client),
// 	Register:   make(chan *Client),
// 	Unregister: make(chan *Client),
// 	Broadcast:  make(chan []byte),
// 	PrivateMsg: make(chan PrivateMessage),
// }

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		// Allow all connections
// 		return true
// 	},
// }

// func (manager *ClientManager) start() {
// 	for {
// 		select {
// 		case client := <-manager.Register:
// 			// Add new client to the map
// 			manager.mutex.Lock()
// 			manager.Clients[client.ID] = client
// 			manager.mutex.Unlock()
// 			fmt.Printf("Client registered: %s (%s)\n", client.Username, client.ID)

// 		case client := <-manager.Unregister:
// 			// Remove client from the map
// 			manager.mutex.Lock()
// 			if _, ok := manager.Clients[client.ID]; ok {
// 				close(client.Send)
// 				delete(manager.Clients, client.ID)
// 				fmt.Printf("Client unregistered: %s (%s)\n", client.Username, client.ID)
// 			}
// 			manager.mutex.Unlock()

// 		case message := <-manager.Broadcast:
// 			// Send message to all clients
// 			manager.mutex.Lock()
// 			for _, client := range manager.Clients {
// 				select {
// 				case client.Send <- message:
// 				default:
// 					close(client.Send)
// 					delete(manager.Clients, client.ID)
// 				}
// 			}
// 			manager.mutex.Unlock()

// 		case privateMsg := <-manager.PrivateMsg:
// 			// Send private message to specific client
// 			manager.mutex.Lock()
// 			for _, client := range manager.Clients {
// 				if client.Username == privateMsg.To {
// 					msg := Message{
// 						Type:    "private",
// 						From:    privateMsg.From,
// 						To:      privateMsg.To,
// 						Content: privateMsg.Content,
// 					}
// 					jsonMsg, _ := json.Marshal(msg)
// 					client.Send <- jsonMsg
// 					fmt.Printf("Private message sent from %s to %s\n", privateMsg.From, privateMsg.To)
// 					break
// 				}
// 			}
// 			manager.mutex.Unlock()
// 		}
// 	}
// }

// func (c *Client) read() {
// 	defer func() {
// 		manager.Unregister <- c
// 		c.Conn.Close()
// 	}()

// 	for {
// 		_, p, err := c.Conn.ReadMessage()
// 		if err != nil {
// 			log.Println("Error reading message:", err)
// 			return
// 		}

// 		// Parse the message
// 		var msg Message
// 		if err := json.Unmarshal(p, &msg); err != nil {
// 			log.Println("Error parsing message:", err)
// 			continue
// 		}

// 		// Handle different message types
// 		switch msg.Type {
// 		case "private":
// 			// Handle private message
// 			manager.PrivateMsg <- PrivateMessage{
// 				From:    c.Username,
// 				To:      msg.To,
// 				Content: msg.Content,
// 			}
// 		case "broadcast":
// 			// Handle broadcast message
// 			manager.Broadcast <- p
// 		case "setUsername":
// 			// Set username for the client
// 			c.Username = msg.Content
// 			fmt.Printf("Username set for client %s: %s\n", c.ID, c.Username)
// 		}
// 	}
// }

// func (c *Client) write() {
// 	defer func() {
// 		c.Conn.Close()
// 	}()

// 	for {
// 		select {
// 		case message, ok := <-c.Send:
// 			if !ok {
// 				// Channel closed
// 				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
// 				return
// 			}

// 			// Write message to the client
// 			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
// 				log.Println("Error writing message:", err)
// 				return
// 			}
// 		}
// 	}
// }

// func handleConnections(w http.ResponseWriter, r *http.Request) {
// 	// Upgrade the HTTP connection to a WebSocket connection
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("Error upgrading connection:", err)
// 		return
// 	}

// 	// Create a new client
// 	client := &Client{
// 		ID:       generateID(),
// 		Conn:     conn,
// 		Send:     make(chan []byte),
// 		Username: "guest",
// 	}

// 	// Register the client
// 	manager.Register <- client

// 	// Start the client's read and write goroutines
// 	go client.read()
// 	go client.write()
// }

// func generateID() string {
// 	return fmt.Sprintf("%d", time.Now().UnixNano())
// }

// func main() {
// 	// Start the client manager
// 	go manager.start()

// 	// Set up the HTTP server
// 	http.HandleFunc("/ws", handleConnections)

// 	// Start the server
// 	fmt.Println("Server starting on :8080")
// 	if err := http.ListenAndServe(":8080", nil); err != nil {
// 		log.Fatal("Error starting server:", err)
// 	}
// }
