package recipe

import (
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/pages"
	"net/http"
	"strconv"
	"strings"
)

func Get(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		success := true
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
		if err != nil {
			success = false
		}
		recipe, err := context.DB.GetRecipeByID(uint(id))
		if err != nil {
			success = false
		}
		if !success {
			messages := []pages.PageMessage{{Type: pages.Error, Value: "Recipe not found"}}
			pages.RenderPage("Recipe", pages.Empty(), messages, w, r)
			return
		}
		directionsParsed := strings.Split(recipe.Directions, "\n")
		pages.RenderPage("Recipe", Recipe(*recipe, directionsParsed), nil, w, r)
	}
}
