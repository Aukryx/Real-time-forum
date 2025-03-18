// web/static/js/user_list.js
import { UserSelectAll } from "./fetch/user.js";
import { initializePrivateMessaging } from "./private_message.js";

// Function to populate the user list on the left of the page
export async function populateUserList() {
  const userList = document.getElementById('userList');
  userList.innerHTML = ''; // Clear existing users
  
  try {
    const userData = await UserSelectAll();
    
    // Add debugging
    // console.log('Fetched user data:', userData);
    
    // Create a section for connected users
    if (userData.connectedUsers && userData.connectedUsers.length > 0) {
      const connectedHeader = document.createElement('li');
      connectedHeader.textContent = 'Connected Users';
      connectedHeader.style.fontWeight = 'bold';
      connectedHeader.style.padding = '0.5rem';
      connectedHeader.style.backgroundColor = '#f0f0f0';
      userList.appendChild(connectedHeader);
      
      userData.connectedUsers.forEach(user => {
        console.log('Processing connected user:', user);
        const li = createUserListItem(user, true);
        userList.appendChild(li);
      });
    }
    
    // Create a section for disconnected users
    if (userData.disconnectedUsers && userData.disconnectedUsers.length > 0) {
      const disconnectedHeader = document.createElement('li');
      disconnectedHeader.textContent = 'Disconnected Users';
      disconnectedHeader.style.fontWeight = 'bold';
      disconnectedHeader.style.padding = '0.5rem';
      disconnectedHeader.style.backgroundColor = '#f0f0f0';
      userList.appendChild(disconnectedHeader);
      
      userData.disconnectedUsers.forEach(user => {
        // console.log('Processing disconnected user:', user);
        const li = createUserListItem(user, false);
        userList.appendChild(li);
      });
    }
    
    // If no users were found in either category
    if ((!userData.connectedUsers || userData.connectedUsers.length === 0) && 
        (!userData.disconnectedUsers || userData.disconnectedUsers.length === 0)) {
      const li = document.createElement('li');
      li.textContent = 'No users found';
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
  
  // Set color based on connection status
  if (isConnected) {
    statusCircle.style.backgroundColor = '#2ecc71'; // Green for connected
  } else {
    statusCircle.style.backgroundColor = '#95a5a6'; // Gray for disconnected
  }
  
  // Username text
  const userText = document.createElement('span');
  userText.textContent = user.nickName || user.NickName || user.name || user.username || `User ${user.id}` || 'Unknown User';
  
  // Add elements to the list item
  li.appendChild(statusCircle);
  li.appendChild(userText);
  
  // We don't need to add the click event here anymore
  // It will be handled by the initializePrivateMessaging function
  
  return li;
}