// websocket.js - Remove the first line and make more targeted exports
export function setupWebSocketForRegistration() {
    const socket = new WebSocket(`ws://${window.location.host}/ws/register`);
    const form = document.getElementById('register-form');
    const messageDiv = document.getElementById('message');
    
    socket.onopen = function() {
        console.log('Connected to registration server');
    };
    
    socket.onclose = function() {
        console.log('Disconnected from server');
        showMessage('Connection to server lost. Please refresh the page.', false);
    };
    
    socket.onerror = function(error) {
        console.error('WebSocket error:', error);
        showMessage('Connection error. Please try again later.', false);
    };
    
    socket.onmessage = function(event) {
        const response = JSON.parse(event.data);
        
        if (response.type === 'register_response') {
            showMessage(response.message, response.success);
            
            if (response.success) {
                // Reset form on success
                form.reset();
                
                // Optional: After successful registration, redirect to login or return to homepage
                setTimeout(() => {
                    const postSection = document.querySelector('.post-section');
                    postSection.innerHTML = '<h2>Posts</h2><p>Registration successful! Welcome to the forum!</p>';
                }, 3000);
            }
        }
    };
    
    if (form) {
        form.addEventListener('submit', function(e) {
            e.preventDefault();
            
            const firstname = document.getElementById('firstname').value;
            const lastname = document.getElementById('lastname').value;
            const username = document.getElementById('username').value;
            const email = document.getElementById('email').value;
            const gender = document.getElementById('gender').value;
            const password = document.getElementById('password').value;
            
            // Send registration request
            socket.send(JSON.stringify({
                type: 'register',
                firstname: firstname,
                lastname: lastname,
                username: username,
                email: email,
                gender: gender,
                password: password
            }));
        });
    }
    
    return socket;
}

function showMessage(text, isSuccess) {
    const messageDiv = document.getElementById('message');
    if (messageDiv) {
        messageDiv.textContent = text;
        messageDiv.className = 'message ' + (isSuccess ? 'success' : 'error');
        messageDiv.style.display = 'block';
    }
}