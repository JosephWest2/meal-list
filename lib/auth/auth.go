package auth

import (
	"errors"
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/db"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

const authCookieName = "meal-list-auth"
const sessionTimeout = 30 * time.Minute

func IsAuthenticated(db *db.DB, r *http.Request) bool {
	authCookie, err := r.Cookie(authCookieName)
	if err != nil {
		return false
	}

	session, err := db.GetSessionByID(authCookie.Value)
	if err != nil {
		return false
	}

	if time.Now().After(session.CreatedAt.Add(sessionTimeout)) {
		db.ClearSession(session.ID)
		return false
	}

	return true
}

func IsAuthorized(db *db.DB, r *http.Request, requiredRole db.Role) bool {
	if !IsAuthenticated(db, r) {
		return false
	}
	session, err := db.GetSessionByID(GetSessionIDFromCookie(r))
	if err != nil {
		return false
	}
	if session.User.Role < requiredRole {
		return false
	}
	return true
}

func Logout(db *db.DB, w http.ResponseWriter, r *http.Request) {
	db.ClearSession(GetSessionIDFromCookie(r))
	ClearSessionCookie(w)
}

// returns true if authentication is successful else false
func Authenticate(db *db.DB, w http.ResponseWriter, r *http.Request, username string, password string) bool {
	user, err := db.GetUserByUsername(username)
	if err != nil {
		// time equalizer for incorrrect username
		bcrypt.CompareHashAndPassword([]byte("jshdf08734t5y467543qwoqu"), []byte(password))
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return false
	}
	sessionId := ulid.Make().String()
	SetSessionCookie(w, sessionId)
	db.CreateSession(sessionId, user.ID)
	return true
}

func GetSessionIDFromCookie(r *http.Request) string {
	authCookie, err := r.Cookie(authCookieName)
	if err != nil {
		return ""
	}
	return authCookie.Value
}

func SetSessionCookie(w http.ResponseWriter, sessionId string) {
	authCookie := http.Cookie{
		Name:     authCookieName,
		Value:    sessionId,
		Path:     "/",
		Expires:  time.Now().Add(30 * time.Minute),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &authCookie)
}

func ClearSessionCookie(w http.ResponseWriter) {
	authCookie := http.Cookie{
		Name:     authCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &authCookie)
}

func GetAuthCookieValue(r *http.Request) string {
	authCookie, err := r.Cookie(authCookieName)
	if err != nil {
		return ""
	}
	return authCookie.Value
}

func WithAuth(requiredRole db.Role, context *app.AppContext, handler http.HandlerFunc) http.HandlerFunc {
	if context == nil {
		panic("context is nil and auth is required")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(context.DB, r) {
			http.Redirect(w, r, "/login?message=Login+Required", http.StatusSeeOther)
			return
		}
		if !IsAuthorized(context.DB, r, requiredRole) {
			http.Redirect(w, r, "/?message=Unauthorized", http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

func IsLoggedInUnverified(r *http.Request) bool {
	return GetAuthCookieValue(r) != ""
}

func GetUserIDFromSession(db *db.DB, r *http.Request) (uint, error) {
	if !IsAuthenticated(db, r) {
		return 0, errors.New("not logged in")
	}
	sessionID := GetSessionIDFromCookie(r)
	if sessionID == "" {
		return 0, errors.New("not logged in")
	}
	session, err := db.GetSessionByID(sessionID)
	if err != nil {
		return 0, errors.New("not logged in")
	}
	return session.UserID, nil
}
