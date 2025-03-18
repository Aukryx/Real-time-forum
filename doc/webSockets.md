# Websockets  

## What are websockets

WebSockets are a communication **protocol** that provides full-duplex communication channels over a single TCP connection.  
Unlike traditional HTTP, which is stateless and requires new connections for each request, WebSockets maintain a persistent connection that allows for **continuous data transfer in both directions.**  
Think of WebSockets like a phone call, where **both parties can speak and listen simultaneously**, rather than traditional HTTP which is more like sending letters back and forth.  

A WebSocket is a persistent, bidirectional communication channel between a single client (a browser/user) and a server. Unlike traditional HTTP requests, which are one-time interactions, WebSockets stay open so that the client and server can send messages anytime.

<br><br>

## Differences between HTTP and WS

### Key Differences Between WebSockets and HTTP Requests

---

| Feature                     | WebSockets (`socket.send()`)       | HTTP POST Request                  |
|-----------------------------|------------------------------------|------------------------------------|
| **Connection**              | Persistent, stays open            | New connection for each request    |
| **Protocol**                | WebSocket (`ws://` or `wss://`)    | HTTP (`http://` or `https://`)     |
| **Overhead**                | Minimal                            | Includes headers, cookies, etc.   |
| **Data Format**             | Raw text, JSON, or binary          | Usually JSON or form data          |
| **Real-Time?**              | Yes, instant message exchange     | No, requires polling or reloading  |

<br>

### Differences  

--- 

| Feature                     | HTTP                                | WebSocket                          |
|-----------------------------|-------------------------------------|------------------------------------|
| **Description**             | A request-response protocol. A client sends a request to a server, which responds with the requested data. | A protocol enabling real-time, bidirectional communication. Once the connection is established, data can flow freely in both directions. |
| **Advantages**              | - Simplicity: Easy to understand and use, widely adopted. <br> - Stateless: Each request is independent, simplifying development. <br> - Large ecosystem: Supported by all browsers and many tools/libraries. <br> - Caching: Responses can be cached, improving performance. | - Real-time communication: Ideal for applications requiring frequent updates (e.g., chats, notifications). <br> - Low latency: Messages can be sent/received without re-establishing a connection. <br> - Efficiency: Less overhead compared to HTTP, as headers are exchanged only once during connection setup. |
| **Disadvantages**           | - Latency: Each request requires a new connection, which can introduce latency. <br> - Not ideal for real-time communication: Not suitable for applications requiring frequent or real-time updates (e.g., online games, chats). <br> - Overhead: Each request and response includes headers, which can lead to significant overhead. | - Complexity: More complex to implement and manage than HTTP. <br> - Connection state: The connection must remain open, which can cause resource management issues. <br> - Limited support: While increasingly supported, some older browsers or environments may not support WebSocket. |
| **Use Case**                | Best for simple applications with sporadic requests. | Best for real-time applications requiring continuous interaction (e.g., chats, live notifications). |

---  

<br><br>
  
# How WebSockets Work  

## 1. Handshake  
The connection begins with a standard HTTP request that includes specific headers indicating a request to upgrade to the WebSocket protocol.  


## 2. Connection Upgrade  
The server acknowledges this request and upgrades the connection from HTTP to WebSocket.  
Once the server accepts the handshake, the connection is established.  

## 3. Data Transfer  
Once established, both client and server can send messages to each other at any time without waiting for a request.  

## 4. Closing  
Either party can close the connection when finished.  

---

<br><br>

# Implementation

## Js - Client side - Init

Usually after the client authenticated on the server side, a response is sent to the client.  
The client will have a websocket connection between him and the server.    

```js
// To get the protocol - http, https, wss or ws
const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
// To get the host - localhost:8080
const host = window.location.host;

const socket = new WebSocket(`${protocol}//${host}/ws`);
```

This **socket** variable contains various informations, such as:
### Properties  

In properties, we have: 
- .url - The full WebSocket URL used to connect (e.g., ws://localhost:8080/ws).  
- .readyState (0 = connecting, 1 = open, 2 = closing, 3 = closed)  
- .protocol - The selected subprotocol (if any) used by the server.  
- .extensions - The WebSocket extensions negotiated with the server.  
- .binaryType - Determines how binary messages are handled ("blob" or "arraybuffer").  

### Methods  

- **.send(data)** to send a message to the server 
- **.close([code, reason])** to close the websocket connection  

### Event Handlers  

- socket.onopen - When the connection is successfully established.  
- socket.onmessage - When a message is received from the server.  
- socket.onerror - When an error occurs (e.g., connection lost).  
- socket.onclose - When the connection is closed.  

## Go - Init  

### Setup the map that will store connected users

```go
var clients = make(map[string]*websocket.Conn)
```

The key is the user's nickName

### Setup the mutex

```go
var mu sync.Mutex
```

The mutex will be there with the Lock() and Unlock() methods to prevent the map from being modified multiple times at the same time.  
We lock, makes the changes and then unlock.  
It's also very important if the map is checked while being modified.  

### Setup the upgrader variable from gorilla  

```go
var upgrader = websocket.Upgrader{
    // Limit the weight of read and write messages
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
```

### Setup the handler  

This handler will be executed as soon as an user is logged, and will have 2 purposes:
- Upgrade the protocol of the user from HTTP to WS
- Listen/Send messages from/to the users

```go
func HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()

    // Storing the client in the connected client map
    // username is retrieved before with cookie or something else
    mu.Lock()
	clients[username] = conn
    mu.Unlock()

    // Infinite loop that will listen to messages
	for {

	}
}
```

### For loop content

We need 2 things in that infinite loop:
- Being constantly reading users messages
- Check if the user is still connected or not

```go
for {
    // Listening to the users messages
    _, msg, err := conn.ReadMessage()
    if err != nil {
        // If error different from nil, that means the connection isn't open anymore
        fmt.Println(username, "disconnected")

        // Removing the user from the connected client map
        mu.Lock()
        delete(clients, username)
        mu.Unlock()
        break
    }
}
```

### Decode the message

We first need to define a struct that will store json content retrieved  

```go
type Message struct {
	Type     string `json:"type"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}
```

We can then decode that json and store it using json.Unmarshal  

```go
var receivedMsg Message
err = json.Unmarshal(msg, &receivedMsg)
if err != nil {
    fmt.Println("Invalid JSON:", err)
    continue
}
```

We then have in that case:
- receivedMsg.Type -> the type of message sent (private message for example)
- receivedMsg.Sender -> the person who sent the message
- receivedMsg.Receiver -> the person supposed to receive the message
- receivedMsg.Message -> the content of the message

The struct and the informations sent via the javascript side depends on the purpose of the websocket connection.  

