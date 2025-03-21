import { getSocket } from "./websockets.js";
import { populateUserList } from "./user_list.js";

// Global variables
let chatWindow = null;
let chatTabs = [];
let currentTab = null;
let currentUsername = null;
let unreadMessages = {}; // Pour suivre les messages non lus par utilisateur

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
    // Ignorer les éléments d'en-tête
    if (item.style.fontWeight === 'bold') return;
    
    // Vérifier si l'utilisateur est connecté
    // const isConnected = item.dataset.connected === 'true';
    
    // Ajouter un événement de clic uniquement pour les utilisateurs connectés
    // if (isConnected) {
      item.addEventListener('click', async (event) => {
        // Obtenir le nom d'utilisateur de l'élément cliqué
        const usernameElement = item.querySelector('span');
        if (!usernameElement) return;
        
        const clickedUsername = usernameElement.textContent;
        
        // Ne pas ouvrir le chat avec soi-même
        if (clickedUsername === currentUsername) {
          console.log('Cannot open chat with yourself');
          return;
        }
        
        // Ouvrir ou créer une fenêtre de chat avec cet utilisateur
        openChatWithUser(clickedUsername);
        
        // Supprimer le point de notification s'il existe
        const notificationDot = item.querySelector('.notification-dot');
        if (notificationDot) {
          notificationDot.remove();
          // Réinitialiser le compteur de messages non lus pour cet utilisateur
          unreadMessages[clickedUsername] = 0;
        }
      });
    // }
  });
}

// Function to handle incoming private messages
export function receivePrivateMessage(sender, messageText) {
  console.log("test");
  
  // Gérer les messages même si la chat window n'est pas créée
  if (!chatWindow) {
    createChatWindow();
    chatWindow.style.display = 'none'; // Créer mais garder caché
  }
  
  // Vérifier si un onglet existe déjà pour cet expéditeur
  const existingTab = chatTabs.find(tab => tab.username === sender);
  let isVisible = false;
  
  // Si la fenêtre de chat existe, qu'elle est visible et que cet onglet est actif
  if (chatWindow && chatWindow.style.display !== 'none' && existingTab && currentTab === existingTab.id) {
    isVisible = true;
  }
  
  // Créer un onglet s'il n'existe pas encore
  if (!existingTab) {
    createNewChatTab(sender);
    // Cacher l'onglet puisqu'on n'affiche pas encore la fenêtre
    if (chatWindow.style.display === 'none') {
      const tabsContainer = document.getElementById('chatTabs');
      if (tabsContainer.lastChild) {
        tabsContainer.lastChild.style.display = 'none';
      }
    }
  }
  
  // Afficher le message reçu dans tous les cas
  displayReceivedMessage(sender, messageText);
  
  // Ajouter la notification si le chat n'est pas visible
  if (!isVisible) {
    addMessageNotification(sender);
  }
  
  // Déplacer l'expéditeur en haut de la liste
  moveUserToTop(sender);
}

// Function to display a received message in the appropriate chat tab
function displayReceivedMessage(sender, messageText) {
  // Create chat window if it doesn't exist
  if (!chatWindow) {
    createChatWindow();
  }
  
  // Check if tab for this user already exists
  let tabData = chatTabs.find(tab => tab.username === sender);
  
  // If no tab exists, create one
  if (!tabData) {
    createNewChatTab(sender);
    tabData = chatTabs.find(tab => tab.username === sender);
  }
  
  if (!tabData) return; // Safety check
  
  // Get the content element
  const contentElement = document.getElementById(tabData.contentId);
  if (!contentElement) return;
  
  // Create a message element
  const messageElement = document.createElement('div');
  messageElement.className = 'message-item received';
  messageElement.style.margin = '5px 0';
  messageElement.style.padding = '8px';
  messageElement.style.borderRadius = '4px';
  messageElement.style.maxWidth = '80%';
  messageElement.style.wordWrap = 'break-word';
  
  // Style as an incoming message
  messageElement.style.backgroundColor = '#e0e0e0';
  messageElement.style.color = 'black';
  messageElement.style.alignSelf = 'flex-start';
  messageElement.style.marginRight = 'auto';
  
  // Add the message text
  messageElement.textContent = messageText;
  
  // Create a message container for flex layout
  const messageContainer = document.createElement('div');
  messageContainer.style.display = 'flex';
  messageContainer.style.flexDirection = 'column';
  messageContainer.appendChild(messageElement);
  
  // Add the message to the content
  contentElement.appendChild(messageContainer);
  
  // Scroll to the bottom
  contentElement.scrollTop = contentElement.scrollHeight;
  
  // Update unread counter
  if (!unreadMessages[sender]) {
    unreadMessages[sender] = 0;
  }
  if (currentTab !== tabData.id) {
    unreadMessages[sender]++;
  }
}

// Function to add notification dot to a user in the list
function addMessageNotification(username) {
  const userListItems = document.querySelectorAll('#userList li');
  
  for (const item of userListItems) {
    const usernameElement = item.querySelector('span');
    if (usernameElement && usernameElement.textContent === username) {
      // Check if notification dot already exists
      let notificationDot = item.querySelector('.notification-dot');
      if (!notificationDot) {
        // Create notification dot
        notificationDot = document.createElement('div');
        notificationDot.className = 'notification-dot';
        notificationDot.style.width = '10px';
        notificationDot.style.height = '10px';
        notificationDot.style.borderRadius = '50%';
        notificationDot.style.backgroundColor = 'red';
        notificationDot.style.marginLeft = 'auto';
        
        // Add it to the user item
        item.appendChild(notificationDot);
      }
      break;
    }
  }
}

// Function to move a user to the top of the user list
function moveUserToTop(username) {
  // // We'll need to call the API to update the user list
  // // For now, we'll just refresh the entire user list
  // // And rely on the server to handle the ordering
  // populateUserList().then(() => {
  //   // After repopulating, re-add any notification dots
  //   const userListItems = document.querySelectorAll('#userList li');
    
  //   for (const item of userListItems) {
  //     const usernameElement = item.querySelector('span');
  //     if (usernameElement && unreadMessages[usernameElement.textContent] > 0) {
  //       // Re-add notification dot
  //       let notificationDot = document.createElement('div');
  //       notificationDot.className = 'notification-dot';
  //       notificationDot.style.width = '10px';
  //       notificationDot.style.height = '10px';
  //       notificationDot.style.borderRadius = '50%';
  //       notificationDot.style.backgroundColor = 'red';
  //       notificationDot.style.marginLeft = 'auto';
        
  //       // Add it to the user item
  //       item.appendChild(notificationDot);
  //     }
  //   }
    
  //   // Re-setup click handlers after DOM changes
  //   setupUserListClickHandlers();
  // });
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
  // [Code existant inchangé pour createChatWindow]
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
  // [Code existant inchangé pour createNewChatTab]
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
  // Masquer tous les contenus d'onglets
  document.querySelectorAll('.tab-content').forEach(content => {
    content.style.display = 'none';
  });
  
  // Enlever la classe active de tous les onglets
  document.querySelectorAll('.chat-tab').forEach(tab => {
    tab.style.backgroundColor = '#f0f0f0';
    tab.style.fontWeight = 'normal';
  });
  
  // Afficher le contenu de l'onglet sélectionné
  const tabContent = document.getElementById(`content_${tabId}`);
  if (tabContent) {
    tabContent.style.display = 'block';
  }
  
  // Mettre la classe active sur l'onglet sélectionné
  const tab = document.getElementById(tabId);
  if (tab) {
    tab.style.backgroundColor = '#e0e0e0';
    tab.style.fontWeight = 'bold';
  }
  
  // Définir l'onglet courant
  currentTab = tabId;
  
  // Trouver le nom d'utilisateur pour cet onglet
  const tabData = chatTabs.find(t => t.id === tabId);
  if (tabData) {
    // Réinitialiser le compteur de messages non lus pour cet utilisateur
    unreadMessages[tabData.username] = 0;
    
    // Supprimer également le point de notification dans la liste des utilisateurs
    const userListItems = document.querySelectorAll('#userList li');
    for (const item of userListItems) {
      const usernameElement = item.querySelector('span');
      if (usernameElement && usernameElement.textContent === tabData.username) {
        const notificationDot = item.querySelector('.notification-dot');
        if (notificationDot) {
          notificationDot.remove();
        }
        break;
      }
    }
  }
  
  // Mettre le focus sur l'input
  const chatInput = document.getElementById('chatInput');
  if (chatInput) {
    chatInput.focus();
  }
}

// Close a specific tab
function closeTab(tabId) {
  // [Code existant inchangé pour closeTab]
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
  // [Code existant inchangé pour sendMessage]
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
}