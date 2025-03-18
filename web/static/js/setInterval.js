import { populatePostList } from "./posts.js";
import { populateUserList } from "./user_list.js";

// Check for new posts every 10 seconds
setInterval(async () => {
    // console.log("Interval");
    
    // try {
    //     const response = await fetch('/api/check-session', {
    //         method: 'GET',
    //         headers: { 'Content-Type': 'application/json' },
    //         // Include credentials to send cookies
    //         credentials: 'include'
    //     });
        
    //     const result = await response.json();
    //     // If the user is logged, aka on the main page
    //     if (result.loggedIn) {
    //         populatePostList()
    //         populateUserList()
    //     }
    // } catch (error) {
    //     console.log("Enable to check session (setInterval.js)");
    // }
}, 10000)