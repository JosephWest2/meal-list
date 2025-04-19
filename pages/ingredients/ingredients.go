package ingredients

import (
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/pages"
	"net/http"
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
		// TODO
	}
}

func Put(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

func Delete(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}
