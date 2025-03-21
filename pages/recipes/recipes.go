package recipes

import (
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/db"
	"josephwest2/meal-list/pages"
	"net/http"
)

func Handler(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			recipes := context.DB.GetRecipes(&db.RecipeQueryParams{})
			pages.RenderPage("Recipes", Recipes(recipes), nil, w, r)
		case "POST":

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
