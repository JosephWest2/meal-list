package router

import (
	"josephwest2/meal-list/api/addRecipeToList"
	"josephwest2/meal-list/api/ingredients"
	"josephwest2/meal-list/api/recipeCategories"
	"josephwest2/meal-list/api/seed"
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"josephwest2/meal-list/lib/sqlc"
	"net/http"
)

func RegisterAPIRoutes(mux *http.ServeMux, app *app.App) {
	mux.HandleFunc("POST /api/seed", seed.Post(app))
	mux.HandleFunc("POST /api/addRecipeToList/{id}", auth.WithAuth(sqlc.RoleStandard, app, addRecipeToList.Post(app)))

	mux.HandleFunc("PATCH /ingredients/{id}", auth.WithAuth(sqlc.RoleAdmin, app, ingredients.Patch(app)))
	mux.HandleFunc("DELETE /ingredients/{id}", auth.WithAuth(sqlc.RoleAdmin, app, ingredients.Delete(app)))

	mux.HandleFunc("PATCH /recipeCategories/{id}", auth.WithAuth(sqlc.RoleAdmin, app, recipeCategories.Patch(app)))
	mux.HandleFunc("DELETE /recipeCategories/{id}", auth.WithAuth(sqlc.RoleAdmin, app, recipeCategories.Delete(app)))
}
