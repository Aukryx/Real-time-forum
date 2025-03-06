// web/static/js/main.js
import { UserSelectAll } from './user.js';
import { fetchPosts, createPost } from './forum.js';

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
  
  const siteTitle = document.createElement('h1');
  siteTitle.textContent = 'PROTS.COM';
  header.appendChild(siteTitle);
  
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
  
  const postInput = document.createElement('input');
  postInput.type = 'text';
  postInput.id = 'newPostInput';
  postInput.placeholder = 'Create a new post...';
  postInput.style.width = '100%';
  postInput.style.padding = '0.5rem';
  postInput.style.marginBottom = '0.5rem';
  
  const submitPostButton = document.createElement('button');
  submitPostButton.textContent = 'Post';
  submitPostButton.id = 'submitPostButton';
  submitPostButton.style.padding = '0.5rem 1rem';
  
  postInputContainer.appendChild(postInput);
  postInputContainer.appendChild(submitPostButton);
  
  const postsTitle = document.createElement('h2');
  postsTitle.textContent = 'Posts';
  
  const postList = document.createElement('ul');
  postList.style.listStyleType = 'none';
  postList.style.padding = '0';
  postList.id = 'postList';
  
  centerColumn.appendChild(postInputContainer);
  centerColumn.appendChild(postsTitle);
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

async function populateUserList() {
    const userList = document.getElementById('userList');
    userList.innerHTML = ''; // Clear existing users
    
    try {
      const users = await UserSelectAll();
      
      // Add debugging
      console.log('Fetched users:', users);
      
      users.forEach(user => {
        console.log('Processing user:', user);
        const li = document.createElement('li');
        li.textContent = user.NickName || user.name || user.username || `User ${user.id}` || 'Unknown User';
        li.style.padding = '0.5rem';
        li.style.borderBottom = '1px solid #ddd';
        li.style.cursor = 'pointer';
        
        li.addEventListener('click', () => {
          console.log('Selected user:', user);
        });
        
        userList.appendChild(li);
      });
    } catch (error) {
      console.error('Error fetching users:', error);
      const li = document.createElement('li');
      li.textContent = 'Error loading users';
      userList.appendChild(li);
    }
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
  const postInput = document.getElementById('newPostInput');
  const submitButton = document.getElementById('submitPostButton');
  
  submitButton.addEventListener('click', async () => {
    const postContent = postInput.value.trim();
    
    if (postContent) {
      try {
        await createPost({
          title: 'New Post', // You might want to add a title input
          content: postContent,
          author: 'CurrentUser' // Replace with actual logged-in user
        });
        
        postInput.value = ''; // Clear input
        await populatePostList(); // Refresh post list
      } catch (error) {
        console.error('Error creating post:', error);
        // Optional: Show error to user
      }
    }
  });
}

// Ensure the page is created when the DOM is loaded
document.addEventListener('DOMContentLoaded', createMainPage);