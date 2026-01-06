package lists

import (
	"josephwest2/meal-list/assert"
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"josephwest2/meal-list/pages"
	"net/http"
)

func Get(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserIDFromSession(context, r)
		assert.Assert(err == nil, "UserID is null in WithAuth protected path")

		listIDString := r.URL.Query().Get("id")
		listID, err := strconv.ParseUint(listIDString, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		list, err := context.Queries.GetListByID(context.QueryContext, int32(listID))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if list.UserID != userID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		associatedRecipes, err := context.Queries.GetAllListAssociatedRecipes(context.QueryContext, int32(listID))
		pages.RenderPage(context, "List", List(&list, associatedRecipes), nil, w, r)
	}
}

