package auth

import (
	"errors"
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/db"
	"josephwest2/meal-list/lib/sqlc"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

const authCookieName = "meal-list-auth"
const sessionTimeout = 24 * time.Hour

const (
	StandardUserRole db.Role = 0
	AdminRole        db.Role = 1
)

func IsAuthenticated(app *app.App, r *http.Request) bool {
	authCookie, err := r.Cookie(authCookieName)
	if err != nil {
		return false
	}

	session, err := app.Queries.GetUserSessionBySessionID(app.QueryContext, authCookie.Value)
	if err != nil {
		return false
	}

	if time.Now().After(session.LastAccess.Time.Add(sessionTimeout)) {
		app.Queries.DeleteUserSessionBySessionID(app.QueryContext, session.ID)
		return false
	}

	return true
}

func IsAuthorized(app *app.App, r *http.Request, requiredRole sqlc.Role) bool {
	if !IsAuthenticated(app, r) {
		return false
	}
	session, err := app.Queries.GetUserSessionBySessionID(app.QueryContext, GetSessionIDFromCookie(r))
	if err != nil {
		return false
	}
	user, err := app.Queries.GetUserByID(app.QueryContext, session.UserID)
	if err != nil {
		return false
	}
	if user.Role < requiredRole {
		return false
	}
	return true
}

func Logout(db *db.DB, w http.ResponseWriter, r *http.Request) {
	db.ClearSession(GetSessionIDFromCookie(r))
	ClearSessionCookie(w)
}

// returns true if authentication is successful else false
func Authenticate(app *app.App, w http.ResponseWriter, r *http.Request, username string, password string) bool {
	user, usernameErr := app.Queries.GetUserByUsername(app.QueryContext, username)
	passErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if usernameErr != nil || passErr != nil {
		return false
	}
	sessionId := ulid.Make().String()
	SetSessionCookie(w, sessionId)
	app.Queries.CreateUserSession(app.QueryContext, sqlc.CreateUserSessionParams{
		ID: sessionId,
		UserID:  user.ID,
	})
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

func WithAuth(requiredRole sqlc.Role, app *app.App, handler http.HandlerFunc) http.HandlerFunc {
	if app == nil {
		panic("app is nil and auth is required")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(app, r) {
			http.Redirect(w, r, "/login?message=Login+Required&target="+r.URL.Path, http.StatusSeeOther)
			return
		}
		if !IsAuthorized(app, r, requiredRole) {
			http.Redirect(w, r, "/?message=Unauthorized", http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

func IsLoggedInUnverified(r *http.Request) bool {
	return GetAuthCookieValue(r) != ""
}

func GetUserIDFromSession(app *app.App, r *http.Request) (int32, error) {
	if !IsAuthenticated(app, r) {
		return 0, errors.New("not logged in")
	}
	sessionID := GetSessionIDFromCookie(r)
	if sessionID == "" {
		return 0, errors.New("not logged in")
	}
	session, err := app.Queries.GetUserSessionBySessionID(app.QueryContext, sessionID)
	if err != nil {
		return 0, errors.New("not logged in")
	}
	return session.UserID, nil
}
