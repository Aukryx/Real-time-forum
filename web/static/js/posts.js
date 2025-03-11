import { fetchPosts } from './fetch/forum.js';

export async function populatePostList() {
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

export function setupPostCreation() {
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