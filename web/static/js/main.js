import { createWelcomePage, removeWelcomePage } from './welcome.js';
import { populateUserList } from './user_list.js';
import { populatePostList, setupPostCreation } from './posts.js';
import { initializePrivateMessaging } from './private_message.js';

// Add this at the top of the file - Demo mode detection
const IS_DEMO_MODE = window.location.hostname.includes('render.com') || 
                    window.location.hostname.includes('onrender.com');

// Add mock WebSocket for demo mode
if (IS_DEMO_MODE) {
  console.log("Running in demo mode - WebSockets disabled");
  
  // Create global mock object to prevent errors
  window.mockWebSocket = {
    send: function() { console.log("WebSocket disabled in demo"); },
    close: function() {}
  };
}

export async function createMainPage() {
  const response = await fetch('/api/navbar', {
    method: 'GET',
    headers: {
        'Content-Type': 'application/json',
    },  
  });

  if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to generate navbar informations: ${errorText}`);
  }

  const user = await response.json();
  // console.log("Response: ", user);
  

  // Set body styles
  document.body.style.fontFamily = 'Arial, sans-serif';
  document.body.style.margin = '0';
  document.body.style.padding = '0';
  document.body.style.display = 'flex';
  document.body.style.flexDirection = 'column';
  document.body.style.minHeight = '100vh';
  document.body.style.backgroundColor = '#f5f5f5';
  
  // Clear body first
  document.body.innerHTML = '';
  
  // Create header
  const header = document.createElement('header');
  header.style.backgroundColor = '#2c3e50';
  header.style.color = 'white';
  header.style.padding = '1rem';
  header.style.textAlign = 'center';
  header.style.display = 'flex';
  header.style.justifyContent = 'space-between';
  header.style.alignItems = 'center';
  
  // Left part of header with site title
  const headerLeft = document.createElement('div');
  const siteTitle = document.createElement('h1');
  siteTitle.textContent = 'PROTS.COM';
  siteTitle.style.margin = '0';
  headerLeft.appendChild(siteTitle);
  
  // Right part of header with button
  const headerRight = document.createElement('div');
  const headerButton = document.createElement('button');
  headerButton.textContent = 'Log Out';
  headerButton.id = 'logout';
  headerButton.style.padding = '0.5rem 1rem';
  headerButton.style.backgroundColor = '#3498db';
  headerButton.style.color = 'white';
  headerButton.style.border = 'none';
  headerButton.style.borderRadius = '4px';
  headerButton.style.cursor = 'pointer';

  const headerName = document.createElement('a');
  headerName.textContent = user.username;
  headerName.id = 'username';
  headerName.style.padding = '0.5rem 3rem';
  headerName.style.color = 'white';
  headerName.style.cursor = 'pointer';
  
  // Add hover effect
  headerButton.addEventListener('mouseenter', () => {
    headerButton.style.backgroundColor = '#2980b9';
  });
  
  headerButton.addEventListener('mouseleave', () => {
    headerButton.style.backgroundColor = '#3498db';
  });
  
  // Add click event listener
  headerButton.addEventListener('click', () => {
    window.location.href = '/logout';
    removeWelcomePage()
    createWelcomePage()
  });
  
  headerRight.appendChild(headerName);
  headerRight.appendChild(headerButton);
  
  // Add both sections to the header
  header.appendChild(headerLeft);
  header.appendChild(headerRight);
  
  // Create content wrapper
  const contentWrapper = document.createElement('div');
  contentWrapper.style.display = 'flex';
  contentWrapper.style.flex = '1';
  contentWrapper.style.position = 'relative';
  
  // Left Column - Users List
  const leftColumn = document.createElement('div');
  leftColumn.style.width = '250px';
  leftColumn.style.backgroundColor = '#f0f0f0';
  leftColumn.style.padding = '1rem';
  leftColumn.style.overflowY = 'auto';
  
  const usersTitle = document.createElement('h2');
  usersTitle.textContent = 'Users';
  leftColumn.appendChild(usersTitle);
  
  const userList = document.createElement('ul');
  userList.style.listStyleType = 'none';
  userList.style.padding = '0';
  userList.id = 'userList';
  leftColumn.appendChild(userList);
  
  // Center Column - Posts
  const centerColumn = document.createElement('div');
  centerColumn.style.flex = '1';
  centerColumn.style.padding = '1rem';
  centerColumn.style.backgroundColor = 'white';
  
  // Post Creation Input
  const postInputContainer = document.createElement('div');
  postInputContainer.style.marginBottom = '1rem';

  const titleInput = document.createElement('input');
  titleInput.type = 'text';
  titleInput.id = 'newPostTitle';
  titleInput.placeholder = 'Post title...';
  titleInput.style.width = '100%';
  titleInput.style.padding = '0.5rem';
  titleInput.style.marginBottom = '0.5rem';
  
  const postInput = document.createElement('input');
  postInput.type = 'text';
  postInput.id = 'newPostInput';
  postInput.placeholder = 'Post content...';
  postInput.style.width = '100%';
  postInput.style.padding = '0.5rem';
  postInput.style.marginBottom = '0.5rem';
  
  const submitPostButton = document.createElement('button');
  submitPostButton.textContent = 'Post';
  submitPostButton.id = 'submitPostButton';
  submitPostButton.style.padding = '0.5rem 1rem';
  
  postInputContainer.appendChild(titleInput);
  postInputContainer.appendChild(postInput);
  postInputContainer.appendChild(submitPostButton);
  
  const postsTitle = document.createElement('h2');
  postsTitle.textContent = 'Posts';
  
  const postList = document.createElement('ul');
  postList.style.listStyleType = 'none';
  postList.style.padding = '0';
  postList.id = 'postList';
  
  const postsContainer = document.createElement('div');
  postsContainer.id = 'posts-container';
  postsContainer.style.marginTop = '1rem';
  postsContainer.style.width = '100%';  // Ensure container has width

  centerColumn.appendChild(postInputContainer);
  centerColumn.appendChild(postsTitle);
  centerColumn.appendChild(postsContainer);
  centerColumn.appendChild(postList);
  
  // Right Column - Random Images
  const rightColumn = document.createElement('div');
  rightColumn.style.width = '250px';
  rightColumn.style.backgroundColor = '#f0f0f0';
  rightColumn.style.padding = '1rem';
  rightColumn.style.overflowY = 'auto';
  
  const imagesTitle = document.createElement('h2');
  imagesTitle.textContent = 'Images';
  rightColumn.appendChild(imagesTitle);
  
  const imageList = document.createElement('ul');
  imageList.style.listStyleType = 'none';
  imageList.style.padding = '0';
  imageList.id = 'imageList';
  rightColumn.appendChild(imageList);
  
  // Add columns to content wrapper
  contentWrapper.appendChild(leftColumn);
  contentWrapper.appendChild(centerColumn);
  contentWrapper.appendChild(rightColumn);
  
  // Footer
  const footer = document.createElement('footer');
  footer.style.backgroundColor = '#2c3e50';
  footer.style.color = 'white';
  footer.style.textAlign = 'center';
  footer.style.padding = '1rem';
  
  const copyright = document.createElement('p');
  copyright.textContent = '© 2025 PROTS.COM. All rights reserved.';
  footer.appendChild(copyright);
  
  // Append to body
  document.body.appendChild(header);
  document.body.appendChild(contentWrapper);
  document.body.appendChild(footer);
  
  // Add demo mode banner if needed
  if (IS_DEMO_MODE) {
    const banner = document.createElement('div');
    banner.style.background = 'rgba(44, 62, 80, 0.9)';
    banner.style.color = 'white';
    banner.style.padding = '10px';
    banner.style.textAlign = 'center';
    banner.style.position = 'fixed';
    banner.style.top = '0';
    banner.style.left = '0';
    banner.style.right = '0';
    banner.style.zIndex = '9999';
    banner.innerHTML = 'DEMO MODE: Real-time chat features are disabled. For the full experience, please check the GitHub repository.';
    document.body.appendChild(banner);
  }
  
  // Initialize everything after DOM elements are created
  initializePage();
  return user.username
}

// New function to handle initialization
function initializePage() {
  populateUserList();
  populatePostList();
  populateImageList();
  setupPostCreation();
  
  // Only initialize messaging if not in demo mode
  if (!IS_DEMO_MODE) {
    initializePrivateMessaging();
  } else {
    console.log("Private messaging disabled in demo mode");
  }
}

function populateImageList() {
  const imageList = document.getElementById('imageList');
  imageList.innerHTML = ''; // Clear existing images
  
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
  
  randomImageIndices.forEach(index => {
    const li = document.createElement('li');
    li.style.marginBottom = '1rem';
    li.style.textAlign = 'center';
    
    const img = document.createElement('img');
    img.src = `/static/img/image-${index}.jpg`;
    img.alt = `Random image ${index}`;
    img.style.maxWidth = '100%';
    img.style.maxHeight = '200px';
    img.style.objectFit = 'cover';
    
    li.appendChild(img);
    imageList.appendChild(li);
  });

}

// Update the DOM loaded event handler
// document.addEventListener('DOMContentLoaded', function() {
//   createMainPage();
// });
