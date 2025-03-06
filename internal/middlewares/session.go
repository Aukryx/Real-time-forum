package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

func GenerateSessionID() string {
	return uuid.New().String()
}

func CreateSession(w http.ResponseWriter, userID int, username, role string) {
	oldSessionID, exists := SessionExists(userID)
	if exists {
		DeleteSession(oldSessionID)
	}
	sessionID := GenerateSessionID()

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true, // Set to true if using HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	// Store session in server (implement this function)
	StoreSession(sessionID, userID, username, role)
}

type Session struct {
	UserID    int
	Username  string
	Role      string
	CreatedAt time.Time
}

var (
	sessions = make(map[string]Session)
	mutex    sync.RWMutex
)

func StoreSession(sessionID string, userID int, username, role string) {
	mutex.Lock()
	defer mutex.Unlock()
	sessions[sessionID] = Session{
		UserID:    userID,
		Username:  username,
		Role:      role,
		CreatedAt: time.Now(),
	}
}

func GetSession(sessionID string) (Session, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	session, exists := sessions[sessionID]
	return session, exists
}

func SessionExists(userID int) (string, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	for sessionID, session := range sessions {
		if session.UserID == userID {
			return sessionID, true
		}
	}
	return "", false
}

func DeleteSession(sessionID string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(sessions, sessionID)
}

func GetCookie(w http.ResponseWriter, r *http.Request) Session {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		//http.Redirect(w, r, "/login", http.StatusSeeOther)
		return Session{}
	}

	sessionID := cookie.Value
	session, exists := GetSession(sessionID)
	if !exists {
		//http.Redirect(w, r, "/login", http.StatusSeeOther)
		return Session{}
	}
	return session
}
