// Function to check if user is logged in
export async function checkSession() {
    try {
      const response = await fetch('/api/check-session', {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' },
        // Include credentials to send cookies
        credentials: 'include'
      });
      
      const result = await response.json();
      return result.loggedIn; // Should return true if session is valid
    } catch (error) {
      console.error('Session check error:', error);
      return false; // Assume not logged in if there's an error
    }
  }