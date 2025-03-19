import { getUsername } from "./getUser.js";
// import { populateUserList } from "./user_list";

let socket = null;

export async function setupWebSockets() {
    let username = await getUsername();
    // Establish WebSocket connection
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host;
    socket = new WebSocket(`${protocol}//${host}/ws`); // Assign to global socket variable

    // Connection opened
    socket.onopen = function () {
        console.log("WebSocket connection established");
        
        // Send a connection message
        // const connectMessage = {
        //     type: 'connect',
        //     username: username
        // };
        // socket.send(JSON.stringify(connectMessage));
    };

    // Handle messages
    socket.onmessage = function (event) {
        try {
            const data = JSON.parse(event.data);
            switch (data.type) {
                case 'private_message':
                    console.log(username, 'received private message:', data);
                    break;
                case 'user_list_update':
                    if (typeof populateUserList === 'function') {
                        populateUserList();
                    }
                    break;
                case 'system_notification':
                    console.log('System notification:', data.message);
                    break;
                default:
                    console.log('Received message:', data);
            }
        } catch (error) {
            console.error('Error parsing WebSocket message:', error);
        }
    };

    // Handle errors
    socket.onerror = function (error) {
        console.error("WebSocket error:", error);
    };

    // Connection closed (attempt to reconnect)
    socket.onclose = function () {
        console.log("WebSocket connection closed. Attempting to reconnect...");
        setTimeout(() => setupWebSockets(username), 3000);
    };

    // Function to send a private message
    socket.sendPrivateMessage = function (receiver, message) {
        console.log(username, "tries to send a private message to", receiver, ":", message);
        
        if (socket.readyState === WebSocket.OPEN) {
            const privateMessage = {
                type: "private_message",
                sender: username,
                receiver: receiver,
                message: message,
            };
            socket.send(JSON.stringify(privateMessage));
        } else {
            console.error("WebSocket is not open. Cannot send message.");
        }
    };

    return socket;
}

// Export a function to get the existing WebSocket instance
export function getSocket() {
    return socket;
}
