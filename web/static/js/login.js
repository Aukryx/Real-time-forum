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

        document.getElementById('login-form').addEventListener('submit', async function (e) {
            e.preventDefault();

            const data = {
                username: document.getElementById('username').value,
                password: document.getElementById('password').value,
            };

            try {
                const response = await fetch('/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(data),
                });

                const result = await response.json();
                showMessage(result.message, result.success);

                if (result.success) {
                    setTimeout(() => {
                        postSection.innerHTML = '<h2>Login successful! Redirecting...</h2>';
                        // Redirect to dashboard or home after login
                        window.location.href = '/dashboard';
                    }, 3000);
                }
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
