package logout

import (
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"josephwest2/meal-list/pages"
	"net/http"
)

func Get(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pages.RenderPage("Logout", Logout(), nil, w, r)
	}
}
func Post(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth.Logout(context.DB, w, r)
		http.Redirect(w, r, "/?message=Logged+Out", http.StatusSeeOther)
	}
}
