// web/static/js/main.js
import { fetchPosts } from './fetch/forum.js';
import { createWelcomePage, removeWelcomePage } from './welcome.js';
import { populateUserList } from './user_list.js';

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
  
  // Initialize everything after DOM elements are created
  initializePage();
}

// New function to handle initialization
function initializePage() {
  populateUserList();
  populatePostList();
  populateImageList();
  setupPostCreation();
  initPostPolling();
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

        const date = new Date(post.CreatedAt);
        console.log("post: ", post);
        
        console.log("date: ", date);
        
        const formattedDate = date.getFullYear() + ' ' + 
          String(date.getMonth() + 1).padStart(2, '0') + ' ' + 
          String(date.getDate()).padStart(2, '0');
        
        const metadata = document.createElement('small');
        metadata.textContent = `By: ${post.Username} | Date: ${formattedDate}`;
        
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

    console.log("Post Content: ", postContent);
    
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
        // console.log('Post created:', newPost);
        
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

// Post submission functionality
const titleInput = document.getElementById('newPostTitle')
const postInput = document.getElementById('newPostInput');
const submitButton = document.getElementById('submitPostButton');
// const postsContainer = document.getElementById('posts-container');

// Function to submit a new post
// Replace your current submitPost function with this one
async function submitPost() {
  const postTitle = titleInput.value.trim();
  const postContent = postInput.value.trim();

  console.log("Post Content: ", postContent);
  
    
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
        // const newPost = await response.json();
        // console.log('Post created:', newPost);
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

// Update the DOM loaded event handler
document.addEventListener('DOMContentLoaded', function() {
  createMainPage();
  // Remove initPostPolling from here since it's now called in createMainPage
});
