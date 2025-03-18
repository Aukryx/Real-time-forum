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

// This should be added to or updated in your forum.js file
export async function fetchPostComments(postId) {
    try {
        const response = await fetch(`/api/posts/${postId}/comments`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Failed to fetch comments: ${errorText}`);
        }

        return await response.json();
    } catch (error) {
        console.error('Error fetching comments:', error);
        throw new Error('Failed to fetch comments');
    }
}