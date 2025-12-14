package recipeCategories

import (
	"josephwest2/meal-list/assert"
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"josephwest2/meal-list/lib/sqlc"
	"josephwest2/meal-list/pages"
	"net/http"
)

func Get(context *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, err := context.Queries.GetAllRecipeCategories(context.QueryContext)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		pages.RenderPage(context, "Recipe Categories", RecipeCategories(categories), nil, w, r)
	}
}

func Post(context *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assert.Assert(auth.IsAuthorized(context, r, sqlc.RoleAdmin), "Post to auth protected route")

	}
}
