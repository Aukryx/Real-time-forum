// web/static/js/welcome.js - Create the welcome page elements dynamically
import { login } from "./login.js";
import { replaceWithRegistrationForm } from "./register.js";

document.addEventListener('DOMContentLoaded', function() {
  createWelcomePage();
});

export function createWelcomePage() {
  // Set body styles first
  document.body.style.fontFamily = 'Arial, sans-serif';
  document.body.style.margin = '0';
  document.body.style.padding = '0';
  document.body.style.display = 'flex';
  document.body.style.flexDirection = 'column';
  document.body.style.minHeight = '100vh';
  document.body.style.backgroundColor = '#f5f5f5';
  
  // Clear body first (in case this function is called multiple times)
  document.body.innerHTML = '';
  
  // Create header
  const header = document.createElement('header');
  header.style.backgroundColor = '#2c3e50';
  header.style.color = 'white';
  header.style.padding = '1rem';
  header.style.textAlign = 'center';
  
  const siteTitle = document.createElement('h1');
  siteTitle.textContent = 'PROTS.COM';
  header.appendChild(siteTitle);
  
  // Create wrapper for content area (between header and footer)
  const contentWrapper = document.createElement('div');
  contentWrapper.style.display = 'flex';
  contentWrapper.style.flex = '1';
  contentWrapper.style.position = 'relative';
  
  // Create left column for fixed images - full height
  const leftColumn = document.createElement('div');
  leftColumn.style.width = '400px';
  leftColumn.style.display = 'flex';
  leftColumn.style.flexDirection = 'column';
  leftColumn.style.backgroundColor = '#f0f0f0';
  
  // Create the 2 left image containers
  for (let i = 0; i < 2; i++) {
    const leftImageDiv = document.createElement('div');
    leftImageDiv.className = 'left-image-container';
    leftImageDiv.style.flex = '1';
    leftImageDiv.style.overflow = 'hidden';
    leftImageDiv.style.position = 'relative';
    
    const leftImage = document.createElement('img');
    leftImage.src = `/static/img/left-image-${i+1}.jpg`; // Fixed images
    leftImage.alt = `Left image ${i+1}`;
    leftImage.style.width = '100%';
    leftImage.style.height = '100%';
    leftImage.style.objectFit = 'cover';
    
    leftImageDiv.appendChild(leftImage);
    leftColumn.appendChild(leftImageDiv);
    
    // Add border between images (except for the last one)
    if (i < 1) {
      const border = document.createElement('div');
      border.style.height = '3px';
      border.style.backgroundColor = 'black';
      leftColumn.appendChild(border);
    }
  }
  
  // Create center column for login
  const centerColumn = document.createElement('div');
  centerColumn.style.flex = '1';
  centerColumn.style.display = 'flex';
  centerColumn.style.justifyContent = 'center';
  centerColumn.style.alignItems = 'center';
  centerColumn.style.padding = '2rem';
  
  // Create login container (rounded square)
  const loginContainer = document.createElement('div');
  loginContainer.className = 'login-container';
  loginContainer.style.backgroundColor = 'white';
  loginContainer.style.borderRadius = '10px';
  loginContainer.style.padding = '2rem';
  loginContainer.style.width = '100%';
  loginContainer.style.maxWidth = '400px';
  loginContainer.style.boxShadow = '0 4px 6px rgba(0, 0, 0, 0.1)';
  
  const welcomeHeading = document.createElement('h2');
  welcomeHeading.textContent = 'Welcome to PROTS.COM';
  loginContainer.appendChild(welcomeHeading);
  
  // Create form
  const loginForm = document.createElement('form');
  loginForm.id = 'loginForm';
  
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
  
  // Button group
  const buttonGroup = document.createElement('div');
  buttonGroup.style.display = 'flex';
  buttonGroup.style.gap = '1rem';
  buttonGroup.style.marginTop = '1.5rem';
  
  // Login button
  const loginButton = document.createElement('button');
  loginButton.type = 'button';
  loginButton.id = 'loginButton';
  loginButton.textContent = 'Log In';
  loginButton.style.flex = '1';
  loginButton.style.padding = '0.75rem';
  loginButton.style.border = 'none';
  loginButton.style.borderRadius = '4px';
  loginButton.style.fontSize = '1rem';
  loginButton.style.cursor = 'pointer';
  loginButton.style.backgroundColor = '#3498db';
  loginButton.style.color = 'white';
  loginButton.addEventListener('mouseover', () => {
    loginButton.style.backgroundColor = '#2980b9';
  });
  loginButton.addEventListener('mouseout', () => {
    loginButton.style.backgroundColor = '#3498db';
  });
  
  // Register button
  const registerButton = document.createElement('button');
  registerButton.type = 'button';
  registerButton.id = 'registerButton';
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
  
  buttonGroup.appendChild(loginButton);
  buttonGroup.appendChild(registerButton);

  let message = document.createElement('div')
  message.setAttribute('class', "message")
  message.setAttribute('style', "display: none;")
  message.setAttribute('id', "message")
  
  // Add all elements to form
  loginForm.appendChild(usernameGroup);
  loginForm.appendChild(passwordGroup);
  loginForm.appendChild(buttonGroup);
  
  // Add form to login container
  loginContainer.appendChild(loginForm);
  loginContainer.appendChild(message)
  
  // Add login container to center column
  centerColumn.appendChild(loginContainer);
  
  // Create right column for random images - full height
  const rightColumn = document.createElement('div');
  rightColumn.style.width = '180px';
  rightColumn.style.display = 'flex';
  rightColumn.style.flexDirection = 'column';
  rightColumn.style.backgroundColor = '#f0f0f0';
  
  // Function to get random images from available pool
  function getRandomImageIndices(totalImages, numToSelect) {
    const indices = [];
    const available = Array.from({ length: totalImages }, (_, i) => i + 1);
    
    for (let i = 0; i < numToSelect; i++) {
      const randomIndex = Math.floor(Math.random() * available.length);
      indices.push(available[randomIndex]);
      available.splice(randomIndex, 1);
    }
    
    return indices;
  }
  
  const totalImages = 17;
  const numImagesToShow = 6;
  const randomImageIndices = getRandomImageIndices(totalImages, numImagesToShow);
  
  // Create the 6 right image containers with random images
  for (let i = 0; i < numImagesToShow; i++) {
    const rightImageDiv = document.createElement('div');
    rightImageDiv.className = 'right-image-container';
    rightImageDiv.style.flex = '1';
    rightImageDiv.style.overflow = 'hidden';
    rightImageDiv.style.position = 'relative';
    rightImageDiv.style.display = 'flex'; // Add flex display
    rightImageDiv.style.justifyContent = 'center'; // Center horizontally
    rightImageDiv.style.alignItems = 'center'; // Center vertically
    rightImageDiv.style.backgroundColor = '#ffffff'; // Light background color
    
    const rightImage = document.createElement('img');
    rightImage.src = `/static/img/image-${randomImageIndices[i]}.jpg`; // Random images
    rightImage.alt = `Random image ${i+1}`;
    rightImage.style.maxWidth = '100%';
    rightImage.style.maxHeight = '100%'; 
    rightImage.style.objectFit = 'contain'; 
    
    rightImageDiv.appendChild(rightImage);
    rightColumn.appendChild(rightImageDiv);
    
    // Add border between images (except for the last one)
    if (i < numImagesToShow - 1) {
      const border = document.createElement('div');
      border.style.height = '3px';
      border.style.backgroundColor = 'black';
      rightColumn.appendChild(border);
    }
  }
  
  // Add all columns to content wrapper
  contentWrapper.appendChild(leftColumn);
  contentWrapper.appendChild(centerColumn);
  contentWrapper.appendChild(rightColumn);
  
  // Create footer
  const footer = document.createElement('footer');
  footer.style.backgroundColor = '#2c3e50';
  footer.style.color = 'white';
  footer.style.textAlign = 'center';
  footer.style.padding = '1rem';
  
  const copyright = document.createElement('p');
  copyright.textContent = '© 2025 PROTS.COM. All rights reserved.';
  footer.appendChild(copyright);
  
  // Add event listeners
  loginButton.addEventListener('click', async function() {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    
    if (username && password) {
      console.log('Login attempt:', { username, password });
      // alert(`Login attempt with username: ${username}`);
      // For a real application, you would handle login logic here
      await login()
    } else {
      // alert('Please enter both username and password');
    }
  });
  
  registerButton.addEventListener('click', function() {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    
    if (username && password) {
      console.log('Registration attempt:', { username, password });

      // alert(`Registration attempt with username: ${username}`);
      // For a real application, you would handle registration logic here
    } else {
      replaceWithRegistrationForm()
      // alert('Please enter both username and password');
    }
  });
  
  // Add all elements to body
  document.body.appendChild(header);
  document.body.appendChild(contentWrapper);
  document.body.appendChild(footer);
}

// Function to remove the welcome page
function removeWelcomePage() {
  document.body.innerHTML = '';
  console.log('Welcome page removed');
}