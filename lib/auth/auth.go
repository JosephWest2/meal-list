package auth

import (
	"errors"
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/sqlc"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

const sessionCookieName = "meal-list-session"
const sessionTimeout = 24 * time.Hour

func IsAuthenticated(context *app.AppContext, r *http.Request) bool {
	sessionCookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return false
	}

	session, err := context.Queries.GetUserSessionBySessionID(context.QueryContext, sessionCookie.Value)
	if err != nil {
		return false
	}

	if time.Now().After(session.LastAccess.Time.Add(sessionTimeout)) {
		context.Queries.DeleteUserSessionBySessionID(context.QueryContext, session.SessionID)
		return false
	}

	return true
}

func IsAuthorized(context *app.AppContext, r *http.Request, requiredRole sqlc.Role) bool {
	if !IsAuthenticated(context, r) {
		return false
	}
	session, err := context.Queries.GetUserSessionBySessionID(context.QueryContext, GetSessionIDFromCookie(r))
	if err != nil {
		return false
	}
	user, err := context.Queries.GetUserByID(context.QueryContext, session.UserID)
	if err != nil {
		return false
	}
	if user.Role < requiredRole {
		return false
	}
	return true
}

func Logout(context *app.AppContext, w http.ResponseWriter, r *http.Request) {
	context.Queries.DeleteUserSessionBySessionID(context.QueryContext, GetSessionIDFromCookie(r))
	ClearSessionCookie(w)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// returns true if authentication is successful else false
func Authenticate(context *app.AppContext, w http.ResponseWriter, r *http.Request, username string, password string) bool {
	user, usernameErr := context.Queries.GetUserByUsername(context.QueryContext, username)
	passErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if usernameErr != nil || passErr != nil {
		return false
	}
	sessionId := ulid.Make().String()
	SetSessionCookie(w, sessionId)
	context.Queries.CreateUserSession(context.QueryContext, sqlc.CreateUserSessionParams{
		SessionID: sessionId,
		UserID:  user.UserID,
	})
	return true
}

func GetSessionIDFromCookie(r *http.Request) string {
	sessionCookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return ""
	}
	return sessionCookie.Value
}

func SetSessionCookie(w http.ResponseWriter, sessionId string) {
	sessionCookie := http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionId,
		Path:     "/",
		Expires:  time.Now().Add(sessionTimeout),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &sessionCookie)
}

func ClearSessionCookie(w http.ResponseWriter) {
	sessionCookie := http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &sessionCookie)
}

func GetSessionCookieValue(r *http.Request) string {
	sessionCookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return ""
	}
	return sessionCookie.Value
}

func WithAuth(requiredRole sqlc.Role, context *app.AppContext, handler http.HandlerFunc) http.HandlerFunc {
	if context == nil {
		panic("app is nil and auth is required")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(context, r) {
			http.Redirect(w, r, "/login?message=Login+Required&redirect-target="+r.URL.Path, http.StatusSeeOther)
			return
		}
		if !IsAuthorized(context, r, requiredRole) {
			http.Redirect(w, r, "/?message=Unauthorized", http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

func IsLoggedInUnverified(r *http.Request) bool {
	return GetSessionCookieValue(r) != ""
}

func GetUserIDFromSession(context *app.AppContext, r *http.Request) (int32, error) {
	if !IsAuthenticated(context, r) {
		return 0, errors.New("not logged in")
	}
	sessionID := GetSessionIDFromCookie(r)
	if sessionID == "" {
		return 0, errors.New("not logged in")
	}
	session, err := context.Queries.GetUserSessionBySessionID(context.QueryContext, sessionID)
	if err != nil {
		return 0, errors.New("not logged in")
	}
	return session.UserID, nil
}
