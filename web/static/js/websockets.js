export function setupWebSockets() {
    // Now establish WebSocket connection
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host;
    const socket = new WebSocket(`${protocol}//${host}/ws`);

    socket.onopen = function(e) {
        console.log("WebSocket connection established");
    };

    socket.onmessage = function(event) {
        console.log("Message received:", event.data);
    };

    socket.onerror = function(error) {
        console.error("WebSocket error:", error);
    };

    socket.onclose = function(event) {
        console.log("WebSocket connection closed", event.code, event.reason);
    };
}