import { getSocket } from "./websockets.js";

// Global variables
let chatWindow = null;
let chatTabs = [];
let currentTab = null;
let currentUsername = null;

// Initialize private messaging functionality
export function initializePrivateMessaging() {
  // Get the current username from the header
  const usernameElement = document.getElementById('username');
  if (usernameElement) {
    currentUsername = usernameElement.textContent;
  }
  
  // Add click event listeners to user list items
  setupUserListClickHandlers();
}

// Set up click handlers for user list
function setupUserListClickHandlers() {
  const userListItems = document.querySelectorAll('#userList li');
  
  userListItems.forEach(item => {
    // Skip the header items
    if (item.style.fontWeight === 'bold') return;
    
    item.addEventListener('click', async (event) => {
      // Get the username from the clicked item
      const usernameElement = item.querySelector('span');
      if (!usernameElement) return;
      
      const clickedUsername = usernameElement.textContent;
      
      // Don't open chat with yourself
      if (clickedUsername === currentUsername) {
        console.log('Cannot open chat with yourself');
        return;
      }
      
      // Open or create chat window with this user
      openChatWithUser(clickedUsername);
    });
  });
}

// Open or create a chat window with a specific user
function openChatWithUser(username) {
  // Create chat window if it doesn't exist
  if (!chatWindow) {
    createChatWindow();
  }
  
  // Check if tab for this user already exists
  const existingTab = chatTabs.find(tab => tab.username === username);
  if (existingTab) {
    // Switch to existing tab
    switchToTab(existingTab.id);
  } else {
    // Create new tab
    createNewChatTab(username);
  }
  
  // Show the chat window
  chatWindow.style.display = 'flex';
}

// Create the chat window
function createChatWindow() {
  // Create the chat window container
  chatWindow = document.createElement('div');
  chatWindow.id = 'chatWindow';
  chatWindow.style.position = 'fixed';
  chatWindow.style.bottom = '0';
  chatWindow.style.right = '20px';
  chatWindow.style.width = '350px';
  chatWindow.style.height = '400px';
  chatWindow.style.backgroundColor = 'white';
  chatWindow.style.border = '1px solid #ccc';
  chatWindow.style.borderRadius = '5px 5px 0 0';
  chatWindow.style.boxShadow = '0 0 10px rgba(0,0,0,0.1)';
  chatWindow.style.display = 'flex';
  chatWindow.style.flexDirection = 'column';
  chatWindow.style.zIndex = '1000';
  
  // Create the header
  const chatHeader = document.createElement('div');
  chatHeader.style.backgroundColor = '#2c3e50';
  chatHeader.style.color = 'white';
  chatHeader.style.padding = '10px';
  chatHeader.style.borderRadius = '5px 5px 0 0';
  chatHeader.style.display = 'flex';
  chatHeader.style.justifyContent = 'space-between';
  chatHeader.style.alignItems = 'center';
  
  // Create the title
  const chatTitle = document.createElement('div');
  chatTitle.textContent = 'Messages';
  chatTitle.style.fontWeight = 'bold';
  
  // Create the close button
  const closeButton = document.createElement('button');
  closeButton.textContent = 'X';
  closeButton.style.backgroundColor = 'transparent';
  closeButton.style.border = 'none';
  closeButton.style.color = 'white';
  closeButton.style.cursor = 'pointer';
  closeButton.style.fontSize = '16px';
  
  // Add click event to close button
  closeButton.addEventListener('click', () => {
    chatWindow.style.display = 'none';
  });
  
  // Add elements to header
  chatHeader.appendChild(chatTitle);
  chatHeader.appendChild(closeButton);
  
  // Create the tabs container
  const tabsContainer = document.createElement('div');
  tabsContainer.id = 'chatTabs';
  tabsContainer.style.display = 'flex';
  tabsContainer.style.backgroundColor = '#f0f0f0';
  tabsContainer.style.borderBottom = '1px solid #ccc';
  tabsContainer.style.overflowX = 'auto';
  
  // Create the chat content area
  const chatContent = document.createElement('div');
  chatContent.id = 'chatContent';
  chatContent.style.flex = '1';
  chatContent.style.overflowY = 'auto';
  chatContent.style.padding = '10px';
  
  // Create the input area
  const inputArea = document.createElement('div');
  inputArea.style.padding = '10px';
  inputArea.style.borderTop = '1px solid #ccc';
  inputArea.style.display = 'flex';
  
  // Create the text input
  const textInput = document.createElement('input');
  textInput.type = 'text';
  textInput.id = 'chatInput';
  textInput.placeholder = 'Type a message...';
  textInput.style.flex = '1';
  textInput.style.padding = '8px';
  textInput.style.border = '1px solid #ccc';
  textInput.style.borderRadius = '4px';
  textInput.style.marginRight = '5px';
  
  // Create the send button
  const sendButton = document.createElement('button');
  sendButton.textContent = 'Send';
  sendButton.style.padding = '8px 12px';
  sendButton.style.backgroundColor = '#3498db';
  sendButton.style.color = 'white';
  sendButton.style.border = 'none';
  sendButton.style.borderRadius = '4px';
  sendButton.style.cursor = 'pointer';
  
  // Add hover effect to send button
  sendButton.addEventListener('mouseenter', () => {
    sendButton.style.backgroundColor = '#2980b9';
  });
  
  sendButton.addEventListener('mouseleave', () => {
    sendButton.style.backgroundColor = '#3498db';
  });
  
  // Add click event to send button
  sendButton.addEventListener('click', () => {
    sendMessage();
  });
  
  // Add keypress event to text input
  textInput.addEventListener('keypress', (event) => {
    if (event.key === 'Enter') {
      sendMessage();
    }
  });
  
  // Add elements to input area
  inputArea.appendChild(textInput);
  inputArea.appendChild(sendButton);
  
  // Add elements to chat window
  chatWindow.appendChild(chatHeader);
  chatWindow.appendChild(tabsContainer);
  chatWindow.appendChild(chatContent);
  chatWindow.appendChild(inputArea);
  
  // Add chat window to document
  document.body.appendChild(chatWindow);
}

// Create a new chat tab
function createNewChatTab(username) {
  // Create a unique ID for the tab
  const tabId = 'tab_' + Date.now();
  
  // Create the tab element
  const tab = document.createElement('div');
  tab.id = tabId;
  tab.className = 'chat-tab';
  tab.style.padding = '8px 15px';
  tab.style.cursor = 'pointer';
  tab.style.position = 'relative';
  tab.style.whiteSpace = 'nowrap';
  tab.style.borderRight = '1px solid #ccc';
  
  // Create the tab name
  const tabName = document.createElement('span');
  tabName.textContent = username;
  
  // Create the close button
  const closeTabButton = document.createElement('span');
  closeTabButton.textContent = '×';
  closeTabButton.style.marginLeft = '8px';
  closeTabButton.style.fontWeight = 'bold';
  
  // Add elements to tab
  tab.appendChild(tabName);
  tab.appendChild(closeTabButton);
  
  // Add click event to tab
  tab.addEventListener('click', (event) => {
    // Only respond to clicks on the tab itself, not the close button
    if (event.target !== closeTabButton) {
      switchToTab(tabId);
    }
  });
  
  // Add click event to close button
  closeTabButton.addEventListener('click', (event) => {
    event.stopPropagation(); // Prevent tab from being selected
    closeTab(tabId);
  });
  
  // Add tab to tabs container
  const tabsContainer = document.getElementById('chatTabs');
  tabsContainer.appendChild(tab);
  
  // Create content for this tab
  const tabContent = document.createElement('div');
  tabContent.id = `content_${tabId}`;
  tabContent.className = 'tab-content';
  tabContent.style.display = 'none';
  tabContent.style.height = '100%';
  tabContent.style.overflowY = 'auto';
  
  // Add welcome message
  const welcomeMessage = document.createElement('div');
  welcomeMessage.textContent = `Start chatting with ${username}`;
  welcomeMessage.style.margin = '10px';
  welcomeMessage.style.color = '#888';
  welcomeMessage.style.textAlign = 'center';
  tabContent.appendChild(welcomeMessage);
  
  // Add content to chat content area
  const chatContent = document.getElementById('chatContent');
  chatContent.appendChild(tabContent);
  
  // Add tab to tabs array
  chatTabs.push({
    id: tabId,
    username: username,
    contentId: `content_${tabId}`
  });
  
  // Switch to the new tab
  switchToTab(tabId);
}

// Switch to a specific tab
function switchToTab(tabId) {
  // Hide all tab contents
  document.querySelectorAll('.tab-content').forEach(content => {
    content.style.display = 'none';
  });
  
  // Remove active class from all tabs
  document.querySelectorAll('.chat-tab').forEach(tab => {
    tab.style.backgroundColor = '#f0f0f0';
    tab.style.fontWeight = 'normal';
  });
  
  // Show the selected tab content
  const tabContent = document.getElementById(`content_${tabId}`);
  if (tabContent) {
    tabContent.style.display = 'block';
  }
  
  // Set active class on the selected tab
  const tab = document.getElementById(tabId);
  if (tab) {
    tab.style.backgroundColor = '#e0e0e0';
    tab.style.fontWeight = 'bold';
  }
  
  // Set currentTab
  currentTab = tabId;
  
  // Focus on the input
  const chatInput = document.getElementById('chatInput');
  if (chatInput) {
    chatInput.focus();
  }
}

// Close a specific tab
function closeTab(tabId) {
  // Find the tab in the array
  const tabIndex = chatTabs.findIndex(tab => tab.id === tabId);
  if (tabIndex === -1) return;
  
  // Remove the tab element
  const tabElement = document.getElementById(tabId);
  if (tabElement) {
    tabElement.remove();
  }
  
  // Remove the tab content
  const contentElement = document.getElementById(`content_${tabId}`);
  if (contentElement) {
    contentElement.remove();
  }
  
  // Remove the tab from the array
  const removedTab = chatTabs.splice(tabIndex, 1)[0];
  
  // If this was the current tab, switch to another tab
  if (currentTab === tabId) {
    if (chatTabs.length > 0) {
      switchToTab(chatTabs[0].id);
    } else {
      // No tabs left, close the window
      chatWindow.style.display = 'none';
      currentTab = null;
    }
  }
}

// Send a message in the current chat
function sendMessage() {
  // Get the input element
  const chatInput = document.getElementById('chatInput');
  if (!chatInput || !currentTab) return;

  // Get the message text
  const messageText = chatInput.value.trim();
  if (!messageText) return;

  // Find the current tab
  const currentTabData = chatTabs.find(tab => tab.id === currentTab);
  if (!currentTabData) return;

  // Get the recipient username
  const receiver = currentTabData.username;
  const socket = getSocket(); // Get the WebSocket instance

  // Get the content element
  const contentElement = document.getElementById(currentTabData.contentId);
  if (!contentElement) return;

  // Create a message element
  const messageElement = document.createElement('div');
  messageElement.className = 'message-item';
  messageElement.style.margin = '5px 0';
  messageElement.style.padding = '8px';
  messageElement.style.borderRadius = '4px';
  messageElement.style.maxWidth = '80%';
  messageElement.style.wordWrap = 'break-word';

  // Style as an outgoing message
  messageElement.style.backgroundColor = '#3498db';
  messageElement.style.color = 'white';
  messageElement.style.alignSelf = 'flex-end';
  messageElement.style.marginLeft = 'auto';

  // Add the message text
  messageElement.textContent = messageText;

  // Create a message container for flex layout
  const messageContainer = document.createElement('div');
  messageContainer.style.display = 'flex';
  messageContainer.style.flexDirection = 'column';
  messageContainer.appendChild(messageElement);

  // Add the message to the content
  contentElement.appendChild(messageContainer);

  // Clear the input
  chatInput.value = '';

  // Scroll to the bottom
  contentElement.scrollTop = contentElement.scrollHeight;

  // Send the message via WebSocket
  if (socket && socket.sendPrivateMessage) {
    socket.sendPrivateMessage(receiver, messageText);
  } else {
    console.error("WebSocket is not initialized or sendPrivateMessage is not defined.");
  }

  // console.log(`Message sent to ${receiver}: ${messageText}`);
}
