// web/static/js/user_list.js
// import {user}
import { initializePrivateMessaging } from "./private_message.js";

// Function to populate the user list on the left of the page
export async function populateUserList(userlist) {
  const userList = document.getElementById('userList');
  userList.innerHTML = ''; // Clear existing users
  
  try {
    // Create a section for connected users
    if (userlist) {
      const connectedHeader = document.createElement('li');
      connectedHeader.textContent = 'Connected Users';
      connectedHeader.style.fontWeight = 'bold';
      connectedHeader.style.padding = '0.5rem';
      connectedHeader.style.backgroundColor = '#f0f0f0';
      userList.appendChild(connectedHeader);
      
      userlist.forEach(user => {
        // console.log('Processing connected user:', user);
        const li = createUserListItem(user, true);
        userList.appendChild(li);
      });
    } else {
      const li = document.createElement('li');
      li.textContent = 'No connected users';
      userList.appendChild(li);
    }
    
    // Initialize private messaging after populating the user list
    initializePrivateMessaging();
  } catch (error) {
    console.error('Error fetching users:', error);
    const li = document.createElement('li');
    li.textContent = 'Error loading users';
    userList.appendChild(li);
  }
}

// Function to create an element for a user in the list
export function createUserListItem(user, isConnected) {
  const li = document.createElement('li');
  
  // Create a flex container for the status circle and username
  li.style.display = 'flex';
  li.style.alignItems = 'center';
  li.style.padding = '0.5rem';
  li.style.borderBottom = '1px solid #ddd';
  li.style.cursor = 'pointer';
  
  // Create status circle
  const statusCircle = document.createElement('div');
  statusCircle.style.width = '12px';
  statusCircle.style.height = '12px';
  statusCircle.style.borderRadius = '50%';
  statusCircle.style.marginRight = '10px';
  
  statusCircle.style.backgroundColor = '#2ecc71';
  
  // Username text
  const userText = document.createElement('span');
  userText.textContent = user;
  
  // Add elements to the list item
  li.appendChild(statusCircle);
  li.appendChild(userText);
  
  return li;
}