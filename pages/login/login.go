package login

import (
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"josephwest2/meal-list/pages"
	"net/http"
)

func Get(context *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pages.RenderPage(context, "Login", Login(), nil, w, r)
	}
}
func Post(context *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")
		target := r.FormValue("redirecttarget")
		if target == "" {
			target = "/"
		}
		if auth.Authenticate(context, w, r, username, password) {
			pages.RedirectWithMessage(w, r, target, pages.PageMessage{Type: pages.MessageSuccess, Value: "Login Success", Timeout: false})
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		messages := []pages.PageMessage{{
			Value:   "Invalid Credentials",
			Type:    pages.MessageError,
			Timeout: false,
		}}
		pages.RenderPage(context, "Login", Login(), messages, w, r)
	}
}
