package router

import (
	"josephwest2/meal-list/api/addRecipeToList"
	"josephwest2/meal-list/api/seed"
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"josephwest2/meal-list/api/ingredients"
	"net/http"
)

func RegisterAPIRoutes(mux *http.ServeMux, context *app.AppContext) {
	mux.HandleFunc("POST /api/seed", seed.Post(context))
	mux.HandleFunc("POST /api/addRecipeToList/{id}", auth.WithAuth(auth.StandardUserRole, context, addRecipeToList.Post(context)))

	mux.HandleFunc("PATCH /ingredients/{id}", ingredients.Patch(context))
	mux.HandleFunc("DELETE /ingredients/{id}", ingredients.Delete(context))
}