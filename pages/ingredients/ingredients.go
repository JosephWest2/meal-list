package ingredients

import (
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/sqlc"
	"josephwest2/meal-list/pages"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

func Get(context *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingredients, err := context.Queries.GetAllIngredients(context.QueryContext)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			println(err.Error())
			return
		}
		categories, err := context.Queries.GetAllIngredientCategories(context.QueryContext)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			println(err.Error())
			return
		}
		pages.RenderPage(context, "Ingredients", Ingredients(ingredients, categories), nil, w, r)
	}
}

func Post(context *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		messages := make([]pages.PageMessage, 0)
		r.ParseForm()
		ingredientName := r.FormValue("ingredient-name")
		categoryIDString := r.FormValue("ingredient-category")
		categoryID, err := strconv.ParseUint(categoryIDString, 10, 32)
		println(ingredientName, categoryIDString)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			println(err.Error())
			return
		}
		_, err = context.Queries.CreateIngredient(context.QueryContext, sqlc.CreateIngredientParams{
			Name:                 ingredientName,
			IngredientCategoryID: pgtype.Int4{Int32: int32(categoryID), Valid: true},
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			println(err.Error())
			return
		}
		messages = append(messages, pages.PageMessage{Type: pages.MessageSuccess, Value: "Ingredient added", Timeout: true})

		ingredients, err := context.Queries.GetAllIngredients(context.QueryContext)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			println(err.Error())
			return
		}
		categories, err := context.Queries.GetAllIngredientCategories(context.QueryContext)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			println(err.Error())
			return
		}
		pages.RenderPage(context, "Ingredients", Ingredients(ingredients, categories), messages, w, r)
	}
}
