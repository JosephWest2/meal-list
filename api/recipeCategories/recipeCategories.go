package recipeCategories

import (
	"josephwest2/meal-list/assert"
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"josephwest2/meal-list/lib/sqlc"
	"net/http"
	"strconv"
)

func Patch(context *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assert.Assert(auth.IsAuthorized(context, r, sqlc.RoleAdmin), "Patch to auth protected route")
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
		_, err = context.Queries.UpdateRecipeCategory(context.QueryContext, sqlc.UpdateRecipeCategoryParams{
			RecipeCategoryID:   int32(id),
			RecipeCategoryName: name})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			print("Failed to update recipe category: " + err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Successfully updated recipe category"))
	}
}

func Delete(context *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assert.Assert(auth.IsAuthorized(context, r, sqlc.RoleAdmin), "Patch to auth protected route")
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = context.Queries.DeleteRecipeCategory(context.QueryContext, int32(id))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			print("Failed to delete recipe category: " + err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Successfully deleted recipe category"))
	}
}
