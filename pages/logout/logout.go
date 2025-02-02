package logout

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
			pages.RenderPage("Logout", Logout(), nil, w, r)
		case "POST":
			auth.Logout(context.DB, w, r)
			http.Redirect(w, r, "/?message=Logged+Out", http.StatusSeeOther)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
