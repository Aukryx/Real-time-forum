package main

import (
	"fmt"
	"log"

	"db"
)

func main() {
	// Initialize database - make sure we're using the actual forum.db
	fmt.Println("Setting up database connection to forum.db...")

	// Initialize the database connection
	database := db.SetupDatabase()
	defer database.Close()

	// Test LikeDislike operations
	fmt.Println("\n==== Testing LikeDislike Operations ====")
	testLikeDislikeOperations()

	// Test Notification operations
	fmt.Println("\n==== Testing Notification Operations ====")
	testNotificationOperations()

	// Add tests for other operations
	fmt.Println("\n==== Testing User Operations ====")
	testUserOperations()

	fmt.Println("\n==== Testing Post Operations ====")
	testPostOperations()

	fmt.Println("\n==== Testing Comment Operations ====")
	testCommentOperations()

	fmt.Println("\n==== Testing Image Operations ====")
	testImageOperations()

	fmt.Println("\n==== Testing Private Message Operations ====")
	testPrivateMessageOperations()

	fmt.Println("\nAll tests completed.")
}

// Existing test functions
func testLikeDislikeOperations() {
	// Use an existing user ID for testing
	userID := 1

	// Test Insert
	postID := 1              // Use an existing post ID
	var commentID *int = nil // No comment ID for this test
	status := 1              // 1 for like, 0 for dislike

	fmt.Println("Creating like/dislike record...")
	likeDislikeID, err := db.LikeDislikeInsert(userID, &postID, commentID, status)
	if err != nil {
		log.Fatalf("Error creating like/dislike: %v\n", err)
	}
	fmt.Printf("Created like/dislike with ID: %d\n", likeDislikeID)

	// Test Select By ID
	fmt.Println("Retrieving like/dislike by ID...")
	likeDislike, err := db.LikeDislikeSelectByID(likeDislikeID)
	if err != nil {
		log.Fatalf("Error retrieving like/dislike: %v\n", err)
	}
	fmt.Printf("Retrieved like/dislike: UserID=%d, PostID=%d, Status=%d\n",
		likeDislike.UserID, likeDislike.PostID, likeDislike.Status)

	// Test Select By Post ID
	fmt.Println("Retrieving likes/dislikes by post ID...")
	likeDislikes, err := db.LikeDislikeSelectByPostID(postID)
	if err != nil {
		log.Fatalf("Error retrieving likes/dislikes: %v\n", err)
	}
	fmt.Printf("Found %d likes/dislikes for post ID %d\n", len(likeDislikes), postID)

	// Test Update
	fmt.Println("Updating like/dislike status...")
	newStatus := 0 // Change to dislike
	err = db.LikeDislikeUpdateStatus(likeDislikeID, newStatus)
	if err != nil {
		log.Fatalf("Error updating like/dislike: %v\n", err)
	}

	// Verify update
	updatedLikeDislike, err := db.LikeDislikeSelectByID(likeDislikeID)
	if err != nil {
		log.Fatalf("Error retrieving updated like/dislike: %v\n", err)
	}
	fmt.Printf("Updated like/dislike status: %d\n", updatedLikeDislike.Status)

	// Keep the record in the database for verification (remove the Delete test)
	fmt.Println("Like/dislike record persisted to database for verification")
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
		log.Fatalf("Error creating notification: %v\n", err)
	}
	fmt.Printf("Created notification with ID: %d\n", notificationID)

	// Test Select By ID
	fmt.Println("Retrieving notification by ID...")
	notification, err := db.NotificationSelectByID(notificationID)
	if err != nil {
		log.Fatalf("Error retrieving notification: %v\n", err)
	}
	fmt.Printf("Retrieved notification: UserID=%d, Type=%s, Content=%s\n",
		notification.UserID, notification.Type, notification.Content)

	// Test Select By User ID
	fmt.Println("Retrieving notifications by user ID...")
	notifications, err := db.NotificationSelectByUserID(int(userID))
	if err != nil {
		log.Fatalf("Error retrieving notifications: %v\n", err)
	}
	fmt.Printf("Found %d notifications for user ID %d\n", len(notifications), userID)

	// Test Update Read Status
	fmt.Println("Updating notification read status...")
	err = db.NotificationUpdateReadStatus(int(notificationID), true)
	if err != nil {
		log.Fatalf("Error updating notification: %v\n", err)
	}

	// Verify update
	updatedNotification, err := db.NotificationSelectByID(int(notificationID))
	if err != nil {
		log.Fatalf("Error retrieving updated notification: %v\n", err)
	}
	fmt.Printf("Updated notification read status: %t\n", updatedNotification.Read)

	// Keep the record in the database for verification (remove the Delete test)
	fmt.Println("Notification record persisted to database for verification")
}

// Additional test functions
func testUserOperations() {
	// This function will test User CRUD operations
	fmt.Println("Creating test user...")
	gender := "male"
	email := "test@example.com"
	firstName := "John"
	lastName := "Doe"
	username := "testuser"
	password := "password123"
	role := "user"

	userID, err := db.UserInsert(username, gender, firstName, lastName, email, password, role)
	if err != "nil" {
		log.Fatalf("Error creating user: %v\n", err)
	}
	fmt.Printf("Created user with ID: %d\n", userID)

	// Test Select By ID
	fmt.Println("Retrieving user by ID...")
	user, errDB := db.UserSelectByID(userID)
	if errDB != nil {
		log.Fatalf("Error retrieving user: %v\n", err)
	}
	fmt.Printf("Retrieved user: Username=%s, Email=%s\n", user.NickName, user.Email)

	// Test Update
	fmt.Println("Updating user information...")
	newUsername := "updateduser"
	errDB = db.UserUpdate(userID, newUsername, gender, firstName, lastName, email, role)
	if errDB != nil {
		log.Fatalf("Error updating user: %v\n", err)
	}

	// Verify update
	updatedUser, errDB := db.UserSelectByID(userID)
	if errDB != nil {
		log.Fatalf("Error retrieving updated user: %v\n", err)
	}
	fmt.Printf("Updated username: %s\n", updatedUser.NickName)

	// Keep the user in the database for other tests to use
}

func testPostOperations() {
	// Use an existing user ID for testing
	userID := 1

	// Test Insert
	title := "Test Post"
	body := "This is a test post content."
	status := "active"
	image := "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTE78YNCrH1vg0c4IoHanLTYYioMGLTvZ8V2w&s"

	fmt.Println("Creating post...")
	postID, err := db.PostInsert(userID, title, body, status)
	if err != nil {
		log.Fatalf("Error creating post: %v\n", err)
	}
	fmt.Printf("Created post with ID: %d\n", postID)

	// Test Select By ID
	fmt.Println("Retrieving post by ID...")
	post, err := db.PostSelectByID(postID)
	if err != nil {
		log.Fatalf("Error retrieving post: %v\n", err)
	}
	fmt.Printf("Retrieved post: Title=%s, UserID=%d\n", post.Title, post.UserID)

	// Test Update
	fmt.Println("Updating post content...")
	newBody := "This is updated post content."
	err = db.PostUpdateContent(postID, title, newBody, image)
	if err != nil {
		log.Fatalf("Error updating post: %v\n", err)
	}

	// Verify update
	updatedPost, err := db.PostSelectByID(postID)
	if err != nil {
		log.Fatalf("Error retrieving updated post: %v\n", err)
	}
	fmt.Printf("Updated post content: %s\n", updatedPost.Body)

	// Keep the post in the database for other tests
}

func testCommentOperations() {
	// Use existing user and post IDs for testing
	userID := 1
	postID := 1

	// Test Insert
	content := "This is a test comment."

	fmt.Println("Creating comment...")
	commentID, err := db.CommentInsert(postID, userID, content)
	if err != nil {
		log.Fatalf("Error creating comment: %v\n", err)
	}
	fmt.Printf("Created comment with ID: %d\n", commentID)

	// Test Select By ID
	fmt.Println("Retrieving comment by ID...")
	comment, err := db.CommentSelectByID(commentID)
	if err != nil {
		log.Fatalf("Error retrieving comment: %v\n", err)
	}
	fmt.Printf("Retrieved comment: Content=%s, UserID=%d\n", comment.Body, comment.UserID)

	// Test Update
	fmt.Println("Updating comment content...")
	newContent := "This is updated comment content."
	err = db.CommentUpdate(commentID, newContent)
	if err != nil {
		log.Fatalf("Error updating comment: %v\n", err)
	}

	// Verify update
	updatedComment, err := db.CommentSelectByID(commentID)
	if err != nil {
		log.Fatalf("Error retrieving updated comment: %v\n", err)
	}
	fmt.Printf("Updated comment content: %s\n", updatedComment.Body)

	// Keep the comment in the database
}

func testImageOperations() {
	// Use an existing post ID for testing
	postID := 1

	// Test Insert
	filePath := "/path/to/test/image.jpg"
	fileSize := 1024 // Size in bytes

	fmt.Println("Creating image record...")
	imageID, err := db.ImageInsert(postID, filePath, fileSize)
	if err != nil {
		log.Fatalf("Error creating image: %v\n", err)
	}
	fmt.Printf("Created image with ID: %d\n", imageID)

	// Test Select By ID
	fmt.Println("Retrieving image by ID...")
	image, err := db.ImageSelectByID(imageID)
	if err != nil {
		log.Fatalf("Error retrieving image: %v\n", err)
	}
	fmt.Printf("Retrieved image: FilePath=%s, PostID=%d\n", image.FilePath, image.PostID)

	// Test Select By Post ID
	fmt.Println("Retrieving images by post ID...")
	images, err := db.ImageSelectByPostID(postID)
	if err != nil {
		log.Fatalf("Error retrieving images: %v\n", err)
	}
	fmt.Printf("Found %d images for post ID %d\n", len(images), postID)

	// Keep the image in the database
}

func testPrivateMessageOperations() {
	// Use existing user IDs for testing
	senderID := 1
	receiverID := 2

	// Test Insert
	message := "This is a test private message."

	fmt.Println("Creating private message...")
	messageID, err := db.PrivateMessageInsert(senderID, receiverID, message)
	if err != nil {
		log.Fatalf("Error creating private message: %v\n", err)
	}
	fmt.Printf("Created private message with ID: %d\n", messageID)

	// Test Select By ID
	fmt.Println("Retrieving private message by ID...")
	privateMessage, err := db.PrivateMessageSelectByID(messageID)
	if err != nil {
		log.Fatalf("Error retrieving private message: %v\n", err)
	}
	fmt.Printf("Retrieved private message: Message=%s, SenderID=%d\n", privateMessage.Message, privateMessage.SenderID)

	// Test Select By User ID
	fmt.Println("Retrieving private messages for user...")
	messages, err := db.PrivateMessageSelectByUserID(receiverID)
	if err != nil {
		log.Fatalf("Error retrieving private messages: %v\n", err)
	}
	fmt.Printf("Found %d private messages for user ID %d\n", len(messages), receiverID)

	// Test Update Read Status
	fmt.Println("Updating private message read status...")
	err = db.PrivateMessageUpdateReadStatus(messageID, true)
	if err != nil {
		log.Fatalf("Error updating private message: %v\n", err)
	}

	// Verify update
	updatedMessage, err := db.PrivateMessageSelectByID(messageID)
	if err != nil {
		log.Fatalf("Error retrieving updated private message: %v\n", err)
	}
	fmt.Printf("Updated private message read status: %t\n", updatedMessage.Read)

	// Keep the private message in the database
}
