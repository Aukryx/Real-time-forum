import { createWelcomePage } from "./welcome.js";
import { removeWelcomePage } from "./welcome.js";
import { createMainPage } from "./main.js";

// Function that creates a div showing the registration result message
export function showMessage(text, isSuccess) {
    const messageDiv = document.getElementById('message');
    if (messageDiv) {
        messageDiv.textContent = text;
        messageDiv.className = 'message ' + (isSuccess ? 'success' : 'error');
        messageDiv.style.display = 'block';
    }
}

// Function to replace login form with registration form
export function replaceWithRegistrationForm() {
    // Get the login container
    const loginContainer = document.querySelector('.login-container');
    if (!loginContainer) return;
    
    // Clear the container
    loginContainer.innerHTML = '';
    
    // Add new heading
    const registerHeading = document.createElement('h2');
    registerHeading.textContent = 'Create an Account';
    loginContainer.appendChild(registerHeading);
    
    // Create registration form
    const registerForm = document.createElement('form');
    registerForm.id = 'registerForm';
    
    // First Name field
    const firstnameGroup = document.createElement('div');
    firstnameGroup.style.marginBottom = '1rem';
    
    const firstnameLabel = document.createElement('label');
    firstnameLabel.textContent = 'First Name';
    firstnameLabel.htmlFor = 'firstname';
    firstnameLabel.style.display = 'block';
    firstnameLabel.style.marginBottom = '0.5rem';
    firstnameLabel.style.fontWeight = 'bold';
    
    const firstnameInput = document.createElement('input');
    firstnameInput.type = 'text';
    firstnameInput.id = 'firstname';
    firstnameInput.name = 'firstname';
    firstnameInput.required = true;
    firstnameInput.style.width = '100%';
    firstnameInput.style.padding = '0.75rem';
    firstnameInput.style.border = '1px solid #ddd';
    firstnameInput.style.borderRadius = '4px';
    firstnameInput.style.fontSize = '1rem';
    firstnameInput.style.boxSizing = 'border-box';
    
    firstnameGroup.appendChild(firstnameLabel);
    firstnameGroup.appendChild(firstnameInput);
    
    // Last Name field
    const lastnameGroup = document.createElement('div');
    lastnameGroup.style.marginBottom = '1rem';
    
    const lastnameLabel = document.createElement('label');
    lastnameLabel.textContent = 'Last Name';
    lastnameLabel.htmlFor = 'lastname';
    lastnameLabel.style.display = 'block';
    lastnameLabel.style.marginBottom = '0.5rem';
    lastnameLabel.style.fontWeight = 'bold';
    
    const lastnameInput = document.createElement('input');
    lastnameInput.type = 'text';
    lastnameInput.id = 'lastname';
    lastnameInput.name = 'lastname';
    lastnameInput.required = true;
    lastnameInput.style.width = '100%';
    lastnameInput.style.padding = '0.75rem';
    lastnameInput.style.border = '1px solid #ddd';
    lastnameInput.style.borderRadius = '4px';
    lastnameInput.style.fontSize = '1rem';
    lastnameInput.style.boxSizing = 'border-box';
    
    lastnameGroup.appendChild(lastnameLabel);
    lastnameGroup.appendChild(lastnameInput);
    
    // Username field
    const usernameGroup = document.createElement('div');
    usernameGroup.style.marginBottom = '1rem';
    
    const usernameLabel = document.createElement('label');
    usernameLabel.textContent = 'Username';
    usernameLabel.htmlFor = 'username';
    usernameLabel.style.display = 'block';
    usernameLabel.style.marginBottom = '0.5rem';
    usernameLabel.style.fontWeight = 'bold';
    
    const usernameInput = document.createElement('input');
    usernameInput.type = 'text';
    usernameInput.id = 'username';
    usernameInput.name = 'username';
    usernameInput.required = true;
    usernameInput.style.width = '100%';
    usernameInput.style.padding = '0.75rem';
    usernameInput.style.border = '1px solid #ddd';
    usernameInput.style.borderRadius = '4px';
    usernameInput.style.fontSize = '1rem';
    usernameInput.style.boxSizing = 'border-box';
    
    usernameGroup.appendChild(usernameLabel);
    usernameGroup.appendChild(usernameInput);
    
    // Email field
    const emailGroup = document.createElement('div');
    emailGroup.style.marginBottom = '1rem';
    
    const emailLabel = document.createElement('label');
    emailLabel.textContent = 'Email';
    emailLabel.htmlFor = 'email';
    emailLabel.style.display = 'block';
    emailLabel.style.marginBottom = '0.5rem';
    emailLabel.style.fontWeight = 'bold';
    
    const emailInput = document.createElement('input');
    emailInput.type = 'email';
    emailInput.id = 'email';
    emailInput.name = 'email';
    emailInput.required = true;
    emailInput.style.width = '100%';
    emailInput.style.padding = '0.75rem';
    emailInput.style.border = '1px solid #ddd';
    emailInput.style.borderRadius = '4px';
    emailInput.style.fontSize = '1rem';
    emailInput.style.boxSizing = 'border-box';
    
    emailGroup.appendChild(emailLabel);
    emailGroup.appendChild(emailInput);
    
    // Gender field
    const genderGroup = document.createElement('div');
    genderGroup.style.marginBottom = '1rem';
    
    const genderLabel = document.createElement('label');
    genderLabel.textContent = 'Gender';
    genderLabel.htmlFor = 'gender';
    genderLabel.style.display = 'block';
    genderLabel.style.marginBottom = '0.5rem';
    genderLabel.style.fontWeight = 'bold';
    
    const genderInput = document.createElement('input');
    genderInput.type = 'text';
    genderInput.id = 'gender';
    genderInput.name = 'gender';
    genderInput.required = true;
    genderInput.style.width = '100%';
    genderInput.style.padding = '0.75rem';
    genderInput.style.border = '1px solid #ddd';
    genderInput.style.borderRadius = '4px';
    genderInput.style.fontSize = '1rem';
    genderInput.style.boxSizing = 'border-box';
    
    genderGroup.appendChild(genderLabel);
    genderGroup.appendChild(genderInput);
    
    // Password field
    const passwordGroup = document.createElement('div');
    passwordGroup.style.marginBottom = '1rem';
    
    const passwordLabel = document.createElement('label');
    passwordLabel.textContent = 'Password';
    passwordLabel.htmlFor = 'password';
    passwordLabel.style.display = 'block';
    passwordLabel.style.marginBottom = '0.5rem';
    passwordLabel.style.fontWeight = 'bold';
    
    const passwordInput = document.createElement('input');
    passwordInput.type = 'password';
    passwordInput.id = 'password';
    passwordInput.name = 'password';
    passwordInput.required = true;
    passwordInput.style.width = '100%';
    passwordInput.style.padding = '0.75rem';
    passwordInput.style.border = '1px solid #ddd';
    passwordInput.style.borderRadius = '4px';
    passwordInput.style.fontSize = '1rem';
    passwordInput.style.boxSizing = 'border-box';
    
    passwordGroup.appendChild(passwordLabel);
    passwordGroup.appendChild(passwordInput);
    
    // Confirm Password field
    const confirmPasswordGroup = document.createElement('div');
    confirmPasswordGroup.style.marginBottom = '1rem';
    
    const confirmPasswordLabel = document.createElement('label');
    confirmPasswordLabel.textContent = 'Confirm Password';
    confirmPasswordLabel.htmlFor = 'confirm-password';
    confirmPasswordLabel.style.display = 'block';
    confirmPasswordLabel.style.marginBottom = '0.5rem';
    confirmPasswordLabel.style.fontWeight = 'bold';
    
    const confirmPasswordInput = document.createElement('input');
    confirmPasswordInput.type = 'password';
    confirmPasswordInput.id = 'confirm-password';
    confirmPasswordInput.name = 'confirm-password';
    confirmPasswordInput.required = true;
    confirmPasswordInput.style.width = '100%';
    confirmPasswordInput.style.padding = '0.75rem';
    confirmPasswordInput.style.border = '1px solid #ddd';
    confirmPasswordInput.style.borderRadius = '4px';
    confirmPasswordInput.style.fontSize = '1rem';
    confirmPasswordInput.style.boxSizing = 'border-box';
    
    confirmPasswordGroup.appendChild(confirmPasswordLabel);
    confirmPasswordGroup.appendChild(confirmPasswordInput);
    
    // Button group
    const buttonGroup = document.createElement('div');
    buttonGroup.style.display = 'flex';
    buttonGroup.style.gap = '1rem';
    buttonGroup.style.marginTop = '1.5rem';
    
    // Register button
    const registerButton = document.createElement('button');
    registerButton.type = 'button';
    registerButton.id = 'registerSubmitButton';
    registerButton.textContent = 'Register';
    registerButton.style.flex = '1';
    registerButton.style.padding = '0.75rem';
    registerButton.style.border = 'none';
    registerButton.style.borderRadius = '4px';
    registerButton.style.fontSize = '1rem';
    registerButton.style.cursor = 'pointer';
    registerButton.style.backgroundColor = '#2ecc71';
    registerButton.style.color = 'white';
    registerButton.addEventListener('mouseover', () => {
      registerButton.style.backgroundColor = '#27ae60';
    });
    registerButton.addEventListener('mouseout', () => {
      registerButton.style.backgroundColor = '#2ecc71';
    });
    
    // Back button
    const backButton = document.createElement('button');
    backButton.type = 'button';
    backButton.id = 'backButton';
    backButton.textContent = 'Back to Login';
    backButton.style.flex = '1';
    backButton.style.padding = '0.75rem';
    backButton.style.border = 'none';
    backButton.style.borderRadius = '4px';
    backButton.style.fontSize = '1rem';
    backButton.style.cursor = 'pointer';
    backButton.style.backgroundColor = '#95a5a6';
    backButton.style.color = 'white';
    backButton.addEventListener('mouseover', () => {
      backButton.style.backgroundColor = '#7f8c8d';
    });
    backButton.addEventListener('mouseout', () => {
      backButton.style.backgroundColor = '#95a5a6';
    });
    
    buttonGroup.appendChild(registerButton);
    buttonGroup.appendChild(backButton);
    
    // Message div for validation/submission feedback
    const messageDiv = document.createElement('div');
    messageDiv.id = 'message';
    messageDiv.className = 'message';
    messageDiv.style.display = 'none';
    messageDiv.style.marginTop = '1rem';
    messageDiv.style.padding = '0.75rem';
    messageDiv.style.borderRadius = '4px';
    
    // Add all elements to form
    registerForm.appendChild(firstnameGroup);
    registerForm.appendChild(lastnameGroup);
    registerForm.appendChild(usernameGroup);
    registerForm.appendChild(emailGroup);
    registerForm.appendChild(genderGroup);
    registerForm.appendChild(passwordGroup);
    registerForm.appendChild(confirmPasswordGroup);
    registerForm.appendChild(buttonGroup);
    
    // Add form to container
    loginContainer.appendChild(registerForm);
    loginContainer.appendChild(messageDiv);
    
    // Event listener for register button
    registerButton.addEventListener('click', async function() {
      const data = {
        firstname: document.getElementById('firstname').value,
        lastname: document.getElementById('lastname').value,
        username: document.getElementById('username').value,
        email: document.getElementById('email').value,
        gender: document.getElementById('gender').value,
        password: document.getElementById('password').value,
        confirm_password: document.getElementById('confirm-password').value
      };
      
      // Check if all fields are filled
      for (const key in data) {
        if (!data[key]) {
          showMessage('Please fill in all fields', false);
          return;
        }
      }
      
      // Check if passwords match
      if (data.password !== data.confirm_password) {
        showMessage("Passwords don't match", false);
        return;
      }
      
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
                // document.querySelector('.post-section').innerHTML = '<h2>Registration successful! Welcome!</h2>';
                removeWelcomePage()
                createMainPage()
            }, 3000);
        }

        // Catching error between the javascript and golang communication
        } catch (error) {
            console.error('Error:', error);
            showMessage('Registration failed. Try again later.', false);
        }
    });
    
    // Event listener for back button
    backButton.addEventListener('click', function() {
      // This would return to the login form
      // In the original code this would call createWelcomePage()
      createWelcomePage();
    });
    
    // Function to show messages
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
  }

  let registerButtontest = document.getElementById('backButton')
  if (registerButtontest) {
    registerButtontest.addEventListener('click', () => {
        createWelcomePage()
      })
  }