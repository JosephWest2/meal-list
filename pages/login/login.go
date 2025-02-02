package login

import (
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"josephwest2/meal-list/pages"
	"net/http"
)

func Handler(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
		case "POST":
			r.ParseForm()
			username := r.FormValue("username")
			password := r.FormValue("password")
			if auth.Authenticate(context.DB, w, r, username, password) {
				pages.RedirectWithMessage(w, r, "/", pages.PageMessage{Type: pages.Success, Value: "Login Success"})
				return
			} else {
				pages.RedirectWithMessage(w, r, "/login", pages.PageMessage{Type: pages.Error, Value: "Invalid Credentials"})
				return
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		pages.RenderPage("Login", Login(), nil, w, r)
	}
}
