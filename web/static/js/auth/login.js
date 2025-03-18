import { createMainPage } from "../main.js";
import { removeWelcomePage } from "../welcome.js";
import { setupWebSockets } from "../websockets.js";

// Get reference to the register and login buttons
// const loginButton = document.querySelector('#loginForm');
// Listener to check when the submit button is pressed
export async function login() {
    // e.preventDefault();
    
    // Getting all the form values and storing them into an object
    const data = {
        username: document.getElementById('username').value,
        password: document.getElementById('password').value,
    };

    // Building the POST request
    try {
        const response = await fetch('/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
        });

        // Await for the response of the golang server and show the message
        const result = await response.json();
        showMessage(result.message, result.success);

        // If the success response if true, send the user to the main page after a short delay
        if (result.success) {
            setTimeout(() => {
                // Redirect to dashboard or home after login
                removeWelcomePage()
                createMainPage()

                setupWebSockets()
            }, 3000);
        }

    // Catching error between the javascript and golang communication
    } catch (error) {
        console.error('Error:', error);
        showMessage('Login failed. Try again later.', false);
    }

    function showMessage(text, isSuccess) {
        const messageDiv = document.getElementById('message');
        if (messageDiv) {
          messageDiv.textContent = text;
          messageDiv.style.display = 'block';
          
          if (isSuccess) {
            messageDiv.style.backgroundColor = '#d4edda';
            messageDiv.style.color = '#155724';
            messageDiv.style.border = '1px solid #c3e6cb';
          } else {
            messageDiv.style.backgroundColor = '#f8d7da';
            messageDiv.style.color = '#721c24';
            messageDiv.style.border = '1px solid #f5c6cb';
          }
        }
      }
};