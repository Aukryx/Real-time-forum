package main

import (
	"db"
	"fmt"
)

func main() {
	// Test database setup
	fmt.Println("Setting up database...")

	// Test LikeDislike operations
	fmt.Println("\n==== Testing LikeDislike Operations ====")
	testLikeDislikeOperations()

	// Test Notification operations
	fmt.Println("\n==== Testing Notification Operations ====")
	testNotificationOperations()

	fmt.Println("\nAll tests completed.")
}

func testLikeDislikeOperations() {
	// Create a test user first (assuming you have a function for this)
	// userID := createTestUser()
	userID := 1 // Or use an existing user ID for testing

	// Test Insert
	postID := 1              // Use an existing post ID
	var commentID *int = nil // No comment ID for this test
	status := 1              // 1 for like, 0 for dislike

	fmt.Println("Creating like/dislike record...")
	likeDislikeID, err := db.LikeDislikeInsert(userID, &postID, commentID, status)
	if err != nil {
		fmt.Printf("Error creating like/dislike: %v\n", err)
		return
	}
	fmt.Printf("Created like/dislike with ID: %d\n", likeDislikeID)

	// Test Select By ID
	fmt.Println("Retrieving like/dislike by ID...")
	likeDislike, err := db.LikeDislikeSelectByID(likeDislikeID)
	if err != nil {
		fmt.Printf("Error retrieving like/dislike: %v\n", err)
		return
	}
	fmt.Printf("Retrieved like/dislike: UserID=%d, PostID=%d, Status=%d\n",
		likeDislike.UserID, likeDislike.PostID, likeDislike.Status)

	// Test Select By Post ID
	fmt.Println("Retrieving likes/dislikes by post ID...")
	likeDislikes, err := db.LikeDislikeSelectByPostID(postID)
	if err != nil {
		fmt.Printf("Error retrieving likes/dislikes: %v\n", err)
		return
	}
	fmt.Printf("Found %d likes/dislikes for post ID %d\n", len(likeDislikes), postID)

	// Test Update
	fmt.Println("Updating like/dislike status...")
	newStatus := 0 // Change to dislike
	err = db.LikeDislikeUpdateStatus(likeDislikeID, newStatus)
	if err != nil {
		fmt.Printf("Error updating like/dislike: %v\n", err)
		return
	}

	// Verify update
	updatedLikeDislike, err := db.LikeDislikeSelectByID(likeDislikeID)
	if err != nil {
		fmt.Printf("Error retrieving updated like/dislike: %v\n", err)
		return
	}
	fmt.Printf("Updated like/dislike status: %d\n", updatedLikeDislike.Status)

	// Test Delete
	fmt.Println("Deleting like/dislike...")
	err = db.LikeDislikeDelete(likeDislikeID)
	if err != nil {
		fmt.Printf("Error deleting like/dislike: %v\n", err)
		return
	}
	fmt.Println("Like/dislike deleted successfully")

	// Verify deletion
	_, err = db.LikeDislikeSelectByID(likeDislikeID)
	if err != nil {
		fmt.Println("Verified: like/dislike no longer exists")
	} else {
		fmt.Println("Error: like/dislike still exists after deletion")
	}
}

func testNotificationOperations() {
	// Use existing user IDs for testing
	userID := 1
	senderID := 2

	// Test Insert
	notificationType := "post_like"
	content := "User 2 liked your post"
	relatedID := 1 // Post ID

	fmt.Println("Creating notification...")
	notificationID, err := db.NotificationInsert(userID, senderID, notificationType, content, relatedID)
	if err != nil {
		fmt.Printf("Error creating notification: %v\n", err)
		return
	}
	fmt.Printf("Created notification with ID: %d\n", notificationID)

	// Test Select By ID
	fmt.Println("Retrieving notification by ID...")
	notification, err := db.NotificationSelectByID(notificationID)
	if err != nil {
		fmt.Printf("Error retrieving notification: %v\n", err)
		return
	}
	fmt.Printf("Retrieved notification: UserID=%d, Type=%s, Content=%s\n",
		notification.UserID, notification.Type, notification.Content)

	// Test Select By User ID
	fmt.Println("Retrieving notifications by user ID...")
	notifications, err := db.NotificationSelectByUserID(userID)
	if err != nil {
		fmt.Printf("Error retrieving notifications: %v\n", err)
		return
	}
	fmt.Printf("Found %d notifications for user ID %d\n", len(notifications), userID)

	// Test Update Read Status
	fmt.Println("Updating notification read status...")
	err = db.NotificationUpdateReadStatus(notificationID, true)
	if err != nil {
		fmt.Printf("Error updating notification: %v\n", err)
		return
	}

	// Verify update
	updatedNotification, err := db.NotificationSelectByID(notificationID)
	if err != nil {
		fmt.Printf("Error retrieving updated notification: %v\n", err)
		return
	}
	fmt.Printf("Updated notification read status: %t\n", updatedNotification.Read)

	// Test Delete
	fmt.Println("Deleting notification...")
	err = db.NotificationDelete(notificationID)
	if err != nil {
		fmt.Printf("Error deleting notification: %v\n", err)
		return
	}
	fmt.Println("Notification deleted successfully")

	// Verify deletion
	_, err = db.NotificationSelectByID(notificationID)
	if err != nil {
		fmt.Println("Verified: notification no longer exists")
	} else {
		fmt.Println("Error: notification still exists after deletion")
	}
}
