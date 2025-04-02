package login

import (
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"josephwest2/meal-list/pages"
	"net/http"
)

func Get(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pages.RenderPage("Login", Login(), nil, w, r)
	}
}
func Post(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")
        target := r.FormValue("redirecttarget")
        if target == "" {
            target = "/"
        }
        println("HELLOO")
		if auth.Authenticate(context.DB, w, r, username, password) {
            println("target: ", target)
			pages.RedirectWithMessage(w, r, target, pages.PageMessage{Type: pages.Success, Value: "Login Success", Timeout: false})
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		messages := []pages.PageMessage{{
			Value:   "Invalid Credentials",
			Type:    pages.Error,
			Timeout: false,
		}}
		pages.RenderPage("Login", Login(), messages, w, r)
	}
}
