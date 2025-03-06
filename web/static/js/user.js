// web/static/js/user.js

export async function UserSelectAll() {
  try {
      const response = await fetch('/api/users', {
          method: 'GET',
          headers: {
              'Content-Type': 'application/json',
          },
      });

      if (!response.ok) {
          const errorText = await response.text();
          throw new Error(`Failed to fetch users: ${errorText}`);
      }

      const users = await response.json();
      console.log('Received users:', users); // Add this debugging line   
      return users;
  } catch (error) {
      console.error('Error fetching users:', error);
      throw error;
  }
}

export async function getUserById(userId) {
  try {
      const response = await fetch(`/api/users/${userId}`, {
          method: 'GET',
          headers: {
              'Content-Type': 'application/json',
          },
      });

      if (!response.ok) {
          throw new Error('Failed to fetch user');
      }

      return await response.json();
  } catch (error) {
      console.error('Error fetching user:', error);
      throw error;
  }
}