package addRecipeToList

import (
	"josephwest2/meal-list/assert"
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"josephwest2/meal-list/lib/sqlc"
	"net/http"
	"strconv"

)

func Post(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assert.Assert(auth.IsAuthenticated(context, r), "Unauthenticated user in WithAuth protected path")

		userID, err := auth.GetUserIDFromSession(context, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			println(err.Error())
			return
		}

		r.ParseForm()

		listIDString := r.FormValue("listID")
		listID, err := strconv.ParseUint(listIDString, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			println(err.Error())
			return
		}

		recipeIDString := r.FormValue("recipeID")
		recipeID, err := strconv.ParseInt(recipeIDString, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			println(err.Error())
			return
		}

		quantityString := r.FormValue("quantity")
		quantity, err := strconv.ParseFloat(quantityString, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			println(err.Error())
			return
		}


		list, err := context.Queries.GetListLeftJoinListToRecipe(context.QueryContext, sqlc.GetListLeftJoinListToRecipeParams{
			ListID: int32(listID),
			RecipeID: int32(recipeID),
		})
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if list.UserID != int32(userID) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if list.ListRecipeQuantity.Valid {
			// lists_to_recipes entry exists, increase quantity
			context.Queries.UpdateListToRecipeQuantity(context.QueryContext, sqlc.UpdateListToRecipeQuantityParams{
				ListRecipeQuantity: float32(quantity) + list.ListRecipeQuantity.Float32,
				ListID: list.ListID,
				RecipeID: int32(recipeID),
			})
		} else {
			// list_to_recipes entry does not exist, create it
			context.Queries.CreateListToRecipe(context.QueryContext, sqlc.CreateListToRecipeParams{
				ListID: list.ListID,
				RecipeID: int32(recipeID),
				ListRecipeQuantity: float32(quantity),
			})
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Successfully added recipe to list"))
	}
}
