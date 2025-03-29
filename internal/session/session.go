package session

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

// Store is the session store
var Store *sessions.CookieStore

// SessionName is the name of the session cookie
const SessionName = "chatp-app-session"

// UserIDKey is the key used to store the user ID in the session
const UserIDKey = "user_id"

// Init initializes the session store
func Init() {
	// Get session secret from environment variable or use a default
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "secret"
	}

	// Create a new cookie store with the secret
	Store = sessions.NewCookieStore([]byte(sessionSecret))

	// Configure the session store
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
	}
}

// SetUserID sets the user ID in the session
func SetUserID(w http.ResponseWriter, r *http.Request, userID int64) error {
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return err
	}
	session.Values[UserIDKey] = userID
	return session.Save(r, w)
}

// GetUserID gets the user ID from the session
func GetUserID(r *http.Request) (int64, bool) {
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return 0, false
	}
	userID, ok := session.Values[UserIDKey].(int64)

	return userID, ok
}

// ClearSession removes all session data
func ClearSession(w http.ResponseWriter, r *http.Request) error {
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return err
	}

	// Clear all session values
	session.Values = make(map[interface{}]interface{})

	// Set MaxAge to -1 to delete the cookie
	session.Options.MaxAge = -1

	return session.Save(r, w)
}

// IsAuthenticated checks if the user is authenticated
func IsAuthenticated(r *http.Request) bool {
	_, ok := GetUserID(r)
	return ok
}
