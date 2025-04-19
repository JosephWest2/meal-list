package recipes

import (
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"josephwest2/meal-list/lib/db"
	"josephwest2/meal-list/pages"
	"net/http"
)

func Get(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbRecipes := context.DB.GetRecipes(&db.RecipeQueryParams{})
		isAdmin := auth.IsAuthorized(context.DB, r, auth.AdminRole)
		pages.RenderPage("Recipes", recipes(dbRecipes, isAdmin), nil, w, r)
	}
}

func Post(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAdmin := auth.IsAuthorized(context.DB, r, auth.AdminRole)
		if !isAdmin {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
