package recipes

import (
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"josephwest2/meal-list/lib/sqlc"
	"josephwest2/meal-list/pages"
	"net/http"
)

func Get(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbRecipes, err := context.Queries.GetRecipesWithOffset(context.QueryContext, sqlc.GetRecipesWithOffsetParams{
			Offset: 0,
			Limit: 50,
		})

		recipeCategoryMap := make(map[int32][]sqlc.RecipeCategory)
		for _, recipe := range dbRecipes {
			categories, err := context.Queries.GetRecipeAssociatedCategories(context.QueryContext, recipe.RecipeID)
			if err != nil {
				continue
			}
			recipeCategoryMap[recipe.RecipeID] = categories
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		isAdmin := auth.IsAuthorized(context, r, sqlc.RoleAdmin)
		pages.RenderPage(context, "Recipes", recipes(dbRecipes, recipeCategoryMap, isAdmin), nil, w, r)
	}
}

func Post(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAdmin := auth.IsAuthorized(context, r, sqlc.RoleAdmin)
		if !isAdmin {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
