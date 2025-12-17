package login

import (
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"josephwest2/meal-list/pages"
	"net/http"
)

func Get(context *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redirectTarget := r.URL.Query().Get("redirect-target")
		pages.RenderPage(context, "Login", Login(redirectTarget), nil, w, r)
	}
}
func Post(context *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")
		redirectTarget := r.URL.Query().Get("redirect-target")
		if redirectTarget == "" {
			redirectTarget = "/"
		}
		if auth.Authenticate(context, w, r, username, password) {
			pages.RedirectWithMessage(w, r, redirectTarget, pages.PageMessage{Type: pages.MessageSuccess, Value: "Login Success"})
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		messages := []pages.PageMessage{{
			Value:   "Invalid Credentials",
			Type:    pages.MessageError,
		}}
		pages.RenderPage(context, "Login", Login(redirectTarget), messages, w, r)
	}
}
