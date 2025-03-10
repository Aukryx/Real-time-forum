// web/static/js/main.js
import { UserSelectAll } from './user.js';
import { fetchPosts, createPost } from './forum.js';
import { createWelcomePage, removeWelcomePage } from './welcome.js';

export function createMainPage() {
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

  centerColumn.appendChild(postInputContainer);
  centerColumn.appendChild(postsTitle);
  centerColumn.appendChild(postsContainer); // Add this line
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
  
  // Populate content
  populateUserList();
  populatePostList();
  populateImageList();
  setupPostCreation();
}

export async function populateUserList() {
  const userList = document.getElementById('userList');
  userList.innerHTML = ''; // Clear existing users
  
  try {
    const userData = await UserSelectAll();
    
    // Add debugging
    console.log('Fetched user data:', userData);
    
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
        console.log('Processing disconnected user:', user);
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
  } catch (error) {
    console.error('Error fetching users:', error);
    const li = document.createElement('li');
    li.textContent = 'Error loading users';
    userList.appendChild(li);
  }
}

function createUserListItem(user, isConnected) {
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
  
  // Add click event to the list item
  li.addEventListener('click', () => {
    console.log('Selected user:', user);
    // You can add additional functionality here, like showing user details
  });
  
  return li;
}

async function populatePostList() {
    const postList = document.getElementById('postList');
    postList.innerHTML = ''; // Clear existing posts
    
    try {
      const posts = await fetchPosts();
      
      // Add debugging
      console.log('Fetched posts:', posts);
      
      if (!posts || posts.length === 0) {
        const li = document.createElement('li');
        li.textContent = 'No posts available';
        postList.appendChild(li);
        return;
      }
  
      posts.forEach(post => {
        const li = document.createElement('li');
        li.style.border = '1px solid #ddd';
        li.style.marginBottom = '1rem';
        li.style.padding = '1rem';
        li.style.borderRadius = '4px';
        
        const title = document.createElement('h3');
        title.textContent = post.Title || 'Untitled Post';
        
        const content = document.createElement('p');
        content.textContent = post.Body || 'No content';
        
        const metadata = document.createElement('small');
        metadata.textContent = `By: User ${post.UserID} | Date: ${post.CreatedAt}`;
        
        li.appendChild(title);
        li.appendChild(content);
        li.appendChild(metadata);
        
        postList.appendChild(li);
      });
    } catch (error) {
      console.error('Error fetching posts:', error);
      const li = document.createElement('li');
      li.textContent = `Error loading posts: ${error.message}`;
      postList.appendChild(li);
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

function setupPostCreation() {
  const titleInput = document.getElementById('newPostTitle');
  const postInput = document.getElementById('newPostInput');
  const submitButton = document.getElementById('submitPostButton');
  
  submitButton.addEventListener('click', async () => {
    const postTitle = titleInput.value.trim();
    const postContent = postInput.value.trim();
    
    if (postTitle && postContent) {
      try {
        const response = await fetch('/api/postCreation', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            title: postTitle,
            content: postContent
          }),
        });
        
        const newPost = await response.json();
        console.log('Post created:', newPost);
        
        // Clear the input fields after successful submission
        titleInput.value = '';
        postInput.value = '';
        
        // Refresh the post list to show the new post
        populatePostList();
      } catch (error) {
        console.error('Error creating post:', error);
      }
    } else {
      console.warn('Post title and content are required');
    }
  });
}

// Throttle function to limit the rate at which we call the API
function throttle(callback, delay) {
  let lastCall = 0;
  return function(...args) {
    const now = new Date().getTime();
    if (now - lastCall < delay) {
      return;
    }
    lastCall = now;
    return callback(...args);
  };
}

// Post submission functionality
const titleInput = document.getElementById('newPostTitle')
const postInput = document.getElementById('newPostInput');
const submitButton = document.getElementById('submitPostButton');
const postsContainer = document.getElementById('posts-container');

// Function to submit a new post
// Replace your current submitPost function with this one
async function submitPost() {
  const postTitle = titleInput.value.trim();
  const postContent = postInput.value.trim();
    
  if (postTitle && postContent) {
    const data = {
      title: postTitle,  // Changed to lowercase to match the setupPostCreation function
      content: postContent  // Changed to 'content' to match the setupPostCreation function
    };

    try {
      const response = await fetch('/api/postCreation', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });
      
      if (!response.ok) {
        const errorText = await response.text();
        console.error(`Server error (${response.status}): ${errorText}`);
        return;
      }
        
      // Try parsing the response as JSON
      try {
        const newPost = await response.json();
        console.log('Post created:', newPost);
      } catch (parseError) {
        console.log('Response was not JSON, but post might have been created');
      }
      
      // Clear the input fields after submission attempt
      titleInput.value = '';
      postInput.value = '';
      
      // Refresh posts
      checkForNewPosts();
      
    } catch (error) {
      console.error('Network error:', error);
    }
  } else {
    console.warn('Post title and content are required');
  }
}

// Add event listener to the submit button
if (submitButton) {
  submitButton.addEventListener('click', submitPost);
}

// Function to fetch new posts from the server
async function fetchNewPosts(lastPostId = 0) {
  try {
    const response = await fetch(`/api/posts/new?lastId=${lastPostId}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    
    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Error fetching new posts:', error);
    return { posts: [] };
  }
}

// Function to render posts to the page
function renderPosts(posts, container) {
  if (!posts || posts.length === 0) return;
  
  posts.forEach(post => {
    const postElement = document.createElement('div');
    postElement.className = 'post';
    postElement.setAttribute('data-id', post.id);
    postElement.innerHTML = `
      <h3>${post.title}</h3>
      <p class="post-meta">Posted by User #${post.user_id} on ${new Date(post.createdAt).toLocaleString()}</p>
      <div class="post-body">${post.body}</div>
      ${post.image ? `<img src="${post.image}" alt="Post image" class="post-image">` : ''}
    `;
    container.prepend(postElement);
  });
}

// Get the largest post ID currently in the DOM
function getLatestPostId() {
  const posts = document.querySelectorAll('.post');
  let maxId = 0;
  
  posts.forEach(post => {
    const postId = parseInt(post.getAttribute('data-id') || '0');
    if (postId > maxId) {
      maxId = postId;
    }
  });
  
  return maxId;
}

// Function to check for new posts
async function checkForNewPosts() {
  if (!postsContainer) return;
  
  const latestPostId = getLatestPostId();
  const { posts } = await fetchNewPosts(latestPostId);
  renderPosts(posts, postsContainer);
}

// Create the throttled function to check for new posts
const checkForNewPostsThrottled = throttle(checkForNewPosts, 5000);

// Initialize post polling
function initPostPolling() {
  if (!postsContainer) {
    console.error("Posts container not found");
    return;
  }
  
  // Initial fetch
  checkForNewPostsThrottled();
  
  // Set up interval for continuous polling
  setInterval(checkForNewPostsThrottled, 5000);
  
  // Also check when the user interacts with the page
  document.addEventListener('visibilitychange', () => {
    if (document.visibilityState === 'visible') {
      checkForNewPostsThrottled();
    }
  });
}

// Initialize when the DOM is fully loaded
document.addEventListener('DOMContentLoaded', function() {
  initPostPolling();
  createMainPage();
});
