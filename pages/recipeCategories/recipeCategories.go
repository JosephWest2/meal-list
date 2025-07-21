package recipeCategories

import (
	"josephwest2/meal-list/assert"
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"josephwest2/meal-list/pages"
	"net/http"
)

func Get(context *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, _ := context.DB.GetAllRecipeCategories()
		pages.RenderPage(context, "Recipe Categories", RecipeCategories(categories), nil, w, r)
	}
}

func Post(context *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assert.Assert(auth.IsAuthorized(context.DB, r, auth.AdminRole), "Post to auth protected route")

	}
}
