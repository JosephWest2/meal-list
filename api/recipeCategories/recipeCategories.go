package recipeCategories

import (
	"josephwest2/meal-list/assert"
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"net/http"
	"strconv"
)

func Patch(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assert.Assert(auth.IsAuthorized(context.DB, r, auth.AdminRole), "Patch to auth protected route")
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			print("Failed to parse form: " + err.Error())
			return
		}
		name := r.FormValue("name")
		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			print("Missing name")
			return
		}
		err = context.DB.UpdateRecipeCategory(uint(id), name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			print("Failed to update recipe category: " + err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Successfully updated recipe category"))
	}
}

func Delete(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assert.Assert(auth.IsAuthorized(context.DB, r, auth.AdminRole), "Patch to auth protected route")
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = context.DB.DeleteRecipeCategory(uint(id))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			print("Failed to delete recipe category: " + err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Successfully deleted recipe category"))
	}
}
