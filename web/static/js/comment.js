import { fetchPostComments } from "./fetch/forum.js";

export async function populateCommentList(postId) {
    const commentList = document.getElementById(`commentList-${postId}`);
    commentList.innerHTML = ''; // Clear existing comments

    try {
        const comments = await fetchPostComments(postId);
        
        // Add debugging
        console.log('Fetched comments:', comments);

        if (!comments || comments.length === 0) {
            const li = document.createElement('li');
            li.textContent = 'No comments available';
            commentList.appendChild(li);
            return;
        }
        comments.forEach(comment => {
            const li = document.createElement('li');
            li.style.border = '1px solid #ddd';
            li.style.marginBottom = '1rem';
            li.style.padding = '1rem';
            li.style.borderRadius = '4px';

            const content = document.createElement('p');
            content.textContent = comment.Body || 'No content';

            const date = new Date(comment.CreatedAt);
            console.log("comment: ", comment);

            console.log("date: ", date);

            const formattedDate = date.getFullYear() + ' ' + 
                String(date.getMonth() + 1).padStart(2, '0') + ' ' + 
                String(date.getDate()).padStart(2, '0');

            const metadata = document.createElement('small');
            metadata.textContent = `By: ${comment.Username} | Date: ${formattedDate}`;

            li.appendChild(content);
            li.appendChild(metadata);

            commentList.appendChild(li);
        });
    } catch (error) {
        console.error('Error fetching comments:', error);
        const li = document.createElement('li');
        li.textContent = `Error loading comments: ${error.message}`;
        commentList.appendChild(li);
    }
}

export function setupCommentCreation(postId) {
    const commentInput = document.getElementById(`newCommentBody-${postId}`);
    const commentButton = document.getElementById(`newCommentButton-${postId}`);

    commentButton.addEventListener('click', async () => {
        const body = commentInput.value.trim();
        if (!body) {
            alert('Please enter a comment');
            return;
        }

        try {
            const response = await fetch(`/api/posts/${postId}/comments`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ body }),
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`Failed to create comment: ${errorText}`);
            }

            commentInput.value = '';
            populateCommentList(postId);
        } catch (error) {
            console.error('Error creating comment:', error);
            alert('Failed to create comment');
        }
    });
}