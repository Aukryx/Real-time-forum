import { showMessage } from "./register.js";

// Get reference to the register and login buttons
const loginButton = document.querySelector('#loginForm');
console.log(loginButton);

// Listener to check when the submit button is pressed
export async function login() {
    // e.preventDefault();
    console.log("test");
    

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
                console.log("success");
                
                let message = document.createElement('div')
                message.innerText = '<h2>Login successful! Redirecting...</h2>'
                document.querySelector('#login-container').appendChild(message)
                // Redirect to dashboard or home after login
                window.location.href = '/dashboard';
            }, 3000);
        }

    // Catching error between the javascript and golang communication
    } catch (error) {
        console.error('Error:', error);
        showMessage('Login failed. Try again later.', false);
    }
};


// Event listener for the login button
if (loginButton) {
    loginButton.addEventListener('click', function (e) {
        e.preventDefault();
        createLoginForm();
    });
}
