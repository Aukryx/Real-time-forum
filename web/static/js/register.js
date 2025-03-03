// Waiting the the page to be fully loaded before doing anything
document.addEventListener('DOMContentLoaded', function () {
    // Get reference to the button that should trigger the registration form
    const registerButton = document.querySelector('.navbar-right button[action="/register"]');
    
    // Get reference to the post section where we'll add the form
    const postSection = document.querySelector('.post-section');

    // Store the original content of the post section
    const originalContent = postSection.innerHTML;

    // Function to create the registration HTML form
    function createRegistrationForm() {
        // Create the registration form HTML
        const formHTML = `
            <h1>Register Account</h1>
            <form id="register-form">
                <div class="form-group">
                    <label for="firstname">First Name:</label>
                    <input type="text" id="firstname" required>
                </div>

                <div class="form-group">
                    <label for="lastname">Last Name:</label>
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
                <div class="form-group">
                    <label for="password">Confirm password:</label>
                    <input type="password" id="confirm-password" required>
                </div>
                
                <button type="submit">Register</button>
            </form>
            <div id="message" class="message" style="display: none;"></div>
        `;

        // Replace the content in the post section
        postSection.innerHTML = formHTML;

        // Listener to check when the submit button is pressed
        document.getElementById('register-form').addEventListener('submit', async function (e) {
            e.preventDefault();

            // Getting all the form values and storing them into an object
            const data = {
                firstname: document.getElementById('firstname').value,
                lastname: document.getElementById('lastname').value,
                username: document.getElementById('username').value,
                email: document.getElementById('email').value,
                gender: document.getElementById('gender').value,
                password: document.getElementById('password').value,
                confirm_password: document.getElementById('confirm-password').value,
            };

            // Checking if the 2 passwords fields match
            if (data.password != data.confirm_password) {
                showMessage("Password doesn't match.", false);
            } else {
                try {
                    // Building the POST request
                    const response = await fetch('/register', {
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
                            document.querySelector('.post-section').innerHTML = '<h2>Registration successful! Welcome!</h2>';
                        }, 3000);
                    }
    
                // Catching error between the javascript and golang communication
                } catch (error) {
                    console.error('Error:', error);
                    showMessage('Registration failed. Try again later.', false);
                }
            }
        });
    }

    // Add click event listener to the register confirmation button
    if (registerButton) {
        registerButton.addEventListener('click', function (e) {
            e.preventDefault();
            createRegistrationForm();
        });
    }

    // Adding the style of the register form
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

// Function that creates a div showing the registration result message
export function showMessage(text, isSuccess) {
    const messageDiv = document.getElementById('message');
    if (messageDiv) {
        messageDiv.textContent = text;
        messageDiv.className = 'message ' + (isSuccess ? 'success' : 'error');
        messageDiv.style.display = 'block';
    }
}
