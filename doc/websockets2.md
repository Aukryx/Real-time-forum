# WebSockets

## What are WebSockets?

WebSockets are a communication **protocol** that provides full-duplex communication channels over a single TCP connection.  
Unlike traditional HTTP, which is stateless and requires new connections for each request, WebSockets maintain a persistent connection that allows for **continuous data transfer in both directions.**  
Think of WebSockets like a phone call, where **both parties can speak and listen simultaneously**, rather than traditional HTTP which is more like sending letters back and forth.  

---

## How WebSockets Work  

### 1. Handshake  
The connection begins with a standard HTTP request that includes specific headers indicating a request to upgrade to the WebSocket protocol.  

---

### 2. Connection Upgrade  
The server acknowledges this request and upgrades the connection from HTTP to WebSocket.  

Once the server accepts the handshake, the connection is established.  

#### **Go Server-side Code:**
```go
package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
    // Upgrading the protocol from http to ws (websocket)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Client connected")
}

func main() {
    // Handler for the websocket innitialisation
	http.HandleFunc("/ws", handleConnection)
	http.ListenAndServe(":8080", nil)
}
```

#### **JavaScript Client-side Code:**
```javascript
const socket = new WebSocket("ws://localhost:8080/ws");

socket.onopen = () => {
    console.log("Connected to WebSocket server");
};
```

---

### 3. Data Transfer  
Once established, both client and server can send messages to each other at any time without waiting for a request.  

#### **Go Server - Sending & Receiving Messages:**
```go
func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
		fmt.Printf("Received: %s\n", msg)

		err = conn.WriteMessage(messageType, []byte("Hello from server!"))
		if err != nil {
			fmt.Println("Write error:", err)
			break
		}
	}
}
```

#### **JavaScript Client - Sending & Receiving Messages:**
```javascript
socket.onmessage = (event) => {
    console.log("Received from server:", event.data);
};

socket.onopen = () => {
    socket.send("Hello from client!");
};
```

---

### 4. Closing  
Either party can close the connection when finished.  

#### **Go Server - Closing the Connection:**
```go
func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Client connected")

	// Example: Close connection after 10 seconds
	time.Sleep(10 * time.Second)
	conn.Close()
}
```

#### **JavaScript Client - Handling Connection Closure:**
```javascript
socket.onclose = () => {
    console.log("WebSocket connection closed");
};
```

---

## Summary
1. **Handshake** - Client requests an upgrade, and the server responds.
2. **Connection Upgrade** - The connection is upgraded to WebSocket.
3. **Data Transfer** - Both client and server can send messages freely.
4. **Closing** - Either party can close the connection at any time.

This allows for real-time communication in web applications, making WebSockets ideal for chat apps, live notifications, and real-time collaboration tools.

## Example  

### One WebSocket Connection Per Client  
Each user connects once to the WebSocket server when they open the application.  
The server keeps track of all connected users.  

Example: Alice and Bob Open the Chat
Alice connects to the WebSocket server.
Bob connects to the WebSocket server.
The server now knows that Alice and Bob are online.
This means that even if Alice is chatting with Bob, and Bob is also chatting with Charlie, they do not need separate WebSocket connections. The server can handle multiple conversations over the same connection.  

### Identifying Who Sends What
Each WebSocket message should include metadata, such as who sent it and who should receive it. The server uses this information to route messages.

Message Format Example (JSON)
When Alice sends Bob a message, it might look like this:

json
Copier
Modifier
{
    "type": "private_message",
    "sender": "Alice",
    "receiver": "Bob",
    "message": "Hey Bob, how are you?"
}
Alice's WebSocket connection sends this JSON message to the server.
The server sees that "receiver": "Bob", so it forwards the message to Bob’s WebSocket connection.




---

## How WebSocket send() Works
Persistent Connection

Unlike HTTP, where each request requires a new connection, WebSockets keep the connection open.
socket.send() sends data through this already open connection, without the overhead of HTTP headers or additional requests.
Message Format

WebSockets work with raw data, meaning they send messages as strings, JSON, or binary data.
In your case, JSON.stringify({...}) converts the message into a JSON string before sending.
What Happens on the Server Side?

The Go WebSocket server listens for incoming messages.
It reads the message, decodes the JSON, and processes it accordingly.  

```js
function sendMessage(receiver, message) {
    socket.send(JSON.stringify({
        type: "private_message",
        sender: "Alice",
        receiver: receiver,
        message: message
    }));
}
```

---
## Readmessage method from websocket.Conn

Understanding ReadMessage() from websocket.Conn
Yes, conn.ReadMessage() is a method provided by the websocket.Conn type from the Gorilla WebSocket package. It reads the next message that arrives on the WebSocket connection.

How ReadMessage() Works:
```go
messageType, msg, err := conn.ReadMessage()
```
- messageType → Specifies the type of message (e.g., websocket.TextMessage or websocket.BinaryMessage).
- msg → Contains the actual message data (in []byte format).
- err → If an error occurs (e.g., client disconnects), it will be handled.
Example Behavior:

If Alice sends a message, her WebSocket connection receives it.
The server reads it using conn.ReadMessage().
But only the server sees this message initially.
→ It does not automatically send the message to other users.

---  

## How Does json.Unmarshal(msg, &receivedMsg) Work?  
When a client (like a user named Alice) sends a WebSocket message, the message is received in byte format.
Since we expect the message to be in JSON format (from socket.send({...}) in JavaScript), we decode it using json.Unmarshal().

Example
If Alice sends:
```json
{
    "type": "private_message",
    "sender": "Alice",
    "receiver": "Bob",
    "message": "Hello, Bob!"
}
```
The Go server receives this as a byte array, then we parse it into a Go struct:  
```go
var receivedMsg Message
err = json.Unmarshal(msg, &receivedMsg)
```
After decoding:
```go
fmt.Println(receivedMsg.Sender)   // Alice
fmt.Println(receivedMsg.Receiver) // Bob
fmt.Println(receivedMsg.Message)  // Hello, Bob!
```
