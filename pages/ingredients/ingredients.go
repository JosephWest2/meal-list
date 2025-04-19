package ingredients

import (
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/pages"
	"net/http"
	"strconv"
)

func Get(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingredients, err := context.DB.GetAllIngredients()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			println(err.Error())
			return
		}
		categories, err := context.DB.GetAllIngredientCategories()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			println(err.Error())
			return
		}
		pages.RenderPage("Ingredients", Ingredients(ingredients, categories), nil, w, r)
	}
}

func Post(context *app.AppContext) http.HandlerFunc {
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
		err = context.DB.CreateIngredient(ingredientName, uint(categoryID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			println(err.Error())
			return
		}
		messages = append(messages, pages.PageMessage{Type: pages.Success, Value: "Ingredient added", Timeout: true})

		ingredients, err := context.DB.GetAllIngredients()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			println(err.Error())
			return
		}
		categories, err := context.DB.GetAllIngredientCategories()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			println(err.Error())
			return
		}
		pages.RenderPage("Ingredients", Ingredients(ingredients, categories), messages, w, r)
	}
}