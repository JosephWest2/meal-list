package recipe

import (
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/pages"
	"net/http"
	"strconv"
	"strings"
)

func Get(context *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		success := true
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
		if err != nil {
			success = false
		}
		recipe, err := context.Queries.GetRecipeByID(context.QueryContext, int32(id))
		if err != nil {
			success = false
		}
		categories, _ := context.Queries.GetRecipeAssociatedCategories(context.QueryContext, int32(id))
		ingredients, _ := context.Queries.GetRecipeAssociatedIngredientsAndUnits(context.QueryContext, int32(id))
		if !success {
			messages := []pages.PageMessage{{Type: pages.MessageError, Value: "Recipe not found"}}
			pages.RenderPage(context, "Recipe", pages.Empty(), messages, w, r)
			return
		}
		directionsParsed := strings.Split(recipe.Directions, "\n")
		pages.RenderPage(context, "Recipe", Recipe(&recipe, categories, ingredients, directionsParsed), nil, w, r)
	}
}
