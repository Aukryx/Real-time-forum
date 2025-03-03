import { showMessage } from "./register.js";

// Waiting for the page to be fully loaded before listening to anything
document.addEventListener('DOMContentLoaded', function () {
    // Get reference to the register and login buttons
    const loginButton = document.querySelector('.navbar-right button[action="/login"]');

    // Get reference to the post section where we'll add the form
    const postSection = document.querySelector('.post-section');

    // Function to create the login form
    function createLoginForm() {
        const formHTML = `
            <h1>Login</h1>
            <form id="login-form">
                <div class="form-group">
                    <label for="username">Username:</label>
                    <input type="text" id="username" required>
                </div>

                <div class="form-group">
                    <label for="password">Password:</label>
                    <input type="password" id="password" required>
                </div>

                <button type="submit">Login</button>
            </form>
            <div id="message" class="message" style="display: none;"></div>
        `;

        postSection.innerHTML = formHTML;

        // Listener to check when the submit button is pressed
        document.getElementById('login-form').addEventListener('submit', async function (e) {
            e.preventDefault();

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
                        postSection.innerHTML = '<h2>Login successful! Redirecting...</h2>';
                        // Redirect to dashboard or home after login
                        window.location.href = '/dashboard';
                    }, 3000);
                }

            // Catching error between the javascript and golang communication
            } catch (error) {
                console.error('Error:', error);
                showMessage('Login failed. Try again later.', false);
            }
        });
    }

    // Event listener for the login button
    if (loginButton) {
        loginButton.addEventListener('click', function (e) {
            e.preventDefault();
            createLoginForm();
        });
    }
});
