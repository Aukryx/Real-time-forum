export async function getUsername() {
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
      return user.username
}