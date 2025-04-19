package ingredients

import (
	"encoding/json"
	"io"
	"josephwest2/meal-list/lib/app"
	"net/http"
	"strconv"
)



type IngredientParams struct {
	Name       string `json:"name"`
	CategoryID uint   `json:"category"`
}

func Patch(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var params IngredientParams
		json.Unmarshal(bodyBytes, &params)
		err = context.DB.UpdateIngredient(uint(id), params.Name, params.CategoryID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			print("Failed to update ingredient: " + err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Successfully updated ingredient"))
	}
}

func Delete(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = context.DB.DeleteIngredient(uint(id))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			print("Failed to delete ingredient: " + err.Error())
			return
		}
	}
}