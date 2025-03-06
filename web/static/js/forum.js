// web/static/js/forum.js

export async function createPost(postData) {
    try {
        const response = await fetch('/api/posts', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                UserID: 1, // Replace with actual user ID
                Title: postData.title || 'Untitled Post',
                Body: postData.content,
                ImagePath: '' // Add image path if needed
            })
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Failed to create post: ${errorText}`);
        }

        return await response.json();
    } catch (error) {
        console.error('Error in createPost:', error);
        throw error;
    }
}

export async function fetchPosts() {
    try {
        const response = await fetch('/api/posts');
        
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Failed to fetch posts: ${errorText}`);
        }

        return await response.json();
    } catch (error) {
        console.error('Error in fetchPosts:', error);
        throw error;
    }
}

export async function fetchPostComments(postId) {
    try {
        const response = await fetch(`/api/posts/${postId}/comments`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            throw new Error('Failed to fetch comments');
        }

        return await response.json();
    } catch (error) {
        console.error('Error fetching comments:', error);
        throw error;
    }
}

export async function createComment(commentData) {
    try {
        const response = await fetch('/api/comments', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                post_id: commentData.postId,
                user_id: getCurrentUserId(), // You'll need to implement this
                body: commentData.content
            }),
        });

        if (!response.ok) {
            throw new Error('Failed to create comment');
        }

        return await response.json();
    } catch (error) {
        console.error('Error creating comment:', error);
        throw error;
    }
}

// Placeholder function - you'll need to implement actual user authentication
function getCurrentUserId() {
    // This should return the ID of the currently logged-in user
    // You might store this in localStorage, sessionStorage, or retrieve from a session
    return 1; // Placeholder
}