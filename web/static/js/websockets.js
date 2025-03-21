import { getUsername } from "./getUser.js";
import { populateUserList } from "./user_list.js";
import { receivePrivateMessage, receiveChatHistory } from "./private_message.js";

let socket = null;

export async function setupWebSockets() {
    // Getting username with a request server
    let username = await getUsername();

    // Establish WebSocket connection
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host;
    socket = new WebSocket(`${protocol}//${host}/ws`); // Assign to global socket variable

    // Method that triggers when the connection is established
    socket.onopen = function () {
        console.log("WebSocket connection established");
    };

    // Method that triggers when an error occurs
    socket.onerror = function (error) {
        console.error("WebSocket error:", error);
    };

    // Method that triggers when connection is closed (attempt to reconnect)
    socket.onclose = function () {
        console.log("WebSocket connection closed. Attempting to reconnect...");
        setTimeout(() => setupWebSockets(username), 3000);
    };

    // Method that triggers when a message is received from the server
    socket.onmessage = function (event) {
        try {
            const data = JSON.parse(event.data);
            // console.log('Received message:', data);
            
            switch (data.type) {
                // When a private message is received
                case 'private_message':
                    receivePrivateMessage(data.sender, data.message);
                    console.log(username, 'received private message:', data);
                    break;
                // When someone connects or disconnects    
                case 'user_list':
                    populateUserList(data.user_list);
                    console.log('User list updated:', data.Userlist);
                    break;
                // When someone opens a chat
                case 'chat_history':
                    console.log('Chat history:', data);
                    receiveChatHistory(data.user2name, data.messages);
                    break;
                case 'system_notification':
                    console.log('System notification:', data.message);
                    break;
                default:
                    console.log('Received data:', data);
            }
        } catch (error) {
            console.error('Error parsing WebSocket message:', error);
        }
    };

    // Function to send a private message
    socket.sendPrivateMessage = function (receiver, message) {
        console.log(username, "Trying to send a private message to", receiver, ":", message);
        
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

    // Function to get the history of messages between 2 users
    socket.getChatHistory = function (receiver) {
        console.log(username, "Requests chat history with", receiver);
        
        // Checking the state of the websocket connection
        if (socket.readyState === WebSocket.OPEN) {
            const chatHistoryRequest = {
                type: "chat_history_request",
                sender: username,
                receiver: receiver,
            };
            socket.send(JSON.stringify(chatHistoryRequest));
        } else {
            console.error("WebSocket is not open. Cannot send message.");
        }
    }

    return socket;
}

// Export a function to get the existing WebSocket instance
export function getSocket() {
    return socket;
}
