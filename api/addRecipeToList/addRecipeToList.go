package addRecipeToList

import (
	"josephwest2/meal-list/assert"
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"net/http"
	"strconv"
)

func Post(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assert.Assert(auth.IsAuthenticated(context.DB, r), "Unauthenticated user in WithAuth protected path")
		userID, err := auth.GetUserIDFromSession(context.DB, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			println(err.Error())
			return
		}
		list := context.DB.GetOrCreateListByUserID(userID)
		recipeIDString := r.PathValue("id")
		recipeID, err := strconv.ParseUint(recipeIDString, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = context.DB.AddRecipeToListOrIncrement(list.ID, uint(recipeID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Successfully added recipe to list"))
	}
}
