package recipes

import (
	"josephwest2/meal-list/lib/db"
	"josephwest2/meal-list/pages"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		if r.FormValue("name") == "seed" {
			db.SeedRecipes()
		}
	} else {
		recipes := db.GetRecipes(&db.RecipeQueryParams{})
		pages.RenderPage("Home", Recipes(recipes), w, r)
	}
}
