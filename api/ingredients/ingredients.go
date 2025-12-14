package ingredients

import (
	"encoding/json"
	"io"
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/sqlc"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

type IngredientParams struct {
	Name       string `json:"name"`
	CategoryID uint   `json:"category"`
}

func Patch(context *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.ParseInt(idString, 10, 32)
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
		_, err = context.Queries.UpdateIngredient(context.QueryContext, sqlc.UpdateIngredientParams{
			ID:                   int32(id),
			Name:                 params.Name,
			IngredientCategoryID: pgtype.Int4{Int32: int32(params.CategoryID), Valid: true},
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			print("Failed to update ingredient: " + err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Successfully updated ingredient"))
	}
}

func Delete(context *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = context.Queries.DeleteIngredient(context.QueryContext, int32(id))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			print("Failed to delete ingredient: " + err.Error())
			return
		}
	}
}
