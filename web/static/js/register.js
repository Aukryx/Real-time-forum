// register.js
import { setupWebSocketForRegistration } from './websocket.js';

document.addEventListener('DOMContentLoaded', function() {
    // Get reference to the button that should trigger the registration form
    const registerButton = document.querySelector('.navbar-right button[action="/ws/register"]');
    
    // Get reference to the post section where we'll add the form
    const postSection = document.querySelector('.post-section');
    
    // Store the original content of the post section
    const originalContent = postSection.innerHTML;
    
    // Create the registration form function
    function createRegistrationForm() {
        // Create the registration form HTML
        const formHTML = `
            <h1>Register Account</h1>
            <form id="register-form">
                <div class="form-group">
                    <label for="first name">First Name:</label>
                    <input type="text" id="firstname" required>
                </div>

                <div class="form-group">
                    <label for="last name">Last Name:</label>
                    <input type="text" id="lastname" required>
                </div>

                <div class="form-group">
                    <label for="username">Username:</label>
                    <input type="text" id="username" required>
                </div>
                
                <div class="form-group">
                    <label for="email">Email:</label>
                    <input type="email" id="email" required>
                </div>

                <div class="form-group">
                    <label for="gender">Gender:</label>
                    <input type="text" id="gender" required>
                </div>
                
                <div class="form-group">
                    <label for="password">Password:</label>
                    <input type="password" id="password" required>
                </div>
                
                <button type="submit">Register</button>
            </form>
            <div id="message" class="message" style="display: none;"></div>
        `;
        
        // Replace the content in the post section
        postSection.innerHTML = formHTML;
        
        // Setup the WebSocket connection
        setupWebSocketForRegistration();
    }
    
    // Add click event listener to the register button
    if (registerButton) {
        registerButton.addEventListener('click', function(e) {
            e.preventDefault();
            createRegistrationForm();
        });
    }
    
    // Add some basic styles if not already in your CSS
    const style = document.createElement('style');
    style.textContent = `
        .form-group { margin-bottom: 15px; }
        label { display: block; margin-bottom: 5px; }
        input { width: 100%; padding: 8px; box-sizing: border-box; }
        .message { margin-top: 20px; padding: 10px; border-radius: 4px; }
        .success { background-color: #dff0d8; color: #3c763d; }
        .error { background-color: #f2dede; color: #a94442; }
    `;
    document.head.appendChild(style);
});