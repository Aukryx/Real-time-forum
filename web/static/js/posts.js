import { fetchPosts } from './fetch/forum.js';
import { populateCommentList, setupCommentCreation } from './comment.js';

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

            // Create comment section
            const commentSection = document.createElement('div');
            commentSection.style.marginTop = '1rem';
            commentSection.style.padding = '1rem';
            commentSection.style.borderTop = '1px solid #ddd';

            const commentTitle = document.createElement('h4');
            commentTitle.textContent = 'Comments';
            commentSection.appendChild(commentTitle);

            const commentList = document.createElement('ul');
            commentList.id = `commentList-${post.ID}`;
            commentList.style.listStyleType = 'none';
            commentList.style.padding = '0';
            commentSection.appendChild(commentList);

            // Add comment input
            const commentInputContainer = document.createElement('div');
            commentInputContainer.style.marginTop = '1rem';

            const commentInput = document.createElement('input');
            commentInput.type = 'text';
            commentInput.id = `newCommentBody-${post.ID}`;
            commentInput.placeholder = 'Add a comment...';
            commentInput.style.width = '100%';
            commentInput.style.padding = '0.5rem';
            commentInput.style.marginBottom = '0.5rem';

            const submitCommentButton = document.createElement('button');
            submitCommentButton.textContent = 'Comment';
            submitCommentButton.id = `newCommentButton-${post.ID}`;
            submitCommentButton.style.padding = '0.5rem 1rem';

            commentInputContainer.appendChild(commentInput);
            commentInputContainer.appendChild(submitCommentButton);
            commentSection.appendChild(commentInputContainer);

            li.appendChild(commentSection);
            postList.appendChild(li);

            // Populate comments and setup comment creation
            populateCommentList(post.ID);
            setupCommentCreation(post.ID);
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