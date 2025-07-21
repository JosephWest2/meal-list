package router

import (
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"josephwest2/meal-list/lib/sqlc"
	"josephwest2/meal-list/pages/index"
	"josephwest2/meal-list/pages/ingredients"
	"josephwest2/meal-list/pages/list"
	"josephwest2/meal-list/pages/login"
	"josephwest2/meal-list/pages/logout"
	recipeCategories "josephwest2/meal-list/pages/recipeCategories"
	"josephwest2/meal-list/pages/recipes"
	createRecipe "josephwest2/meal-list/pages/recipes/createRecipe"
	"josephwest2/meal-list/pages/recipes/recipe"
	"josephwest2/meal-list/pages/register"
	"net/http"
)

func RegisterPageRoutes(mux *http.ServeMux, context *app.App) {

	mux.HandleFunc("GET /", index.Get(context))

	mux.HandleFunc("GET /recipes", recipes.Get(context))
	mux.HandleFunc("POST /recipes", auth.WithAuth(sqlc.RoleAdmin, context, recipes.Post(context)))

	mux.HandleFunc("GET /recipes/create", auth.WithAuth(sqlc.RoleAdmin, context, createRecipe.Get(context)))
	mux.HandleFunc("POST /recipes/create", auth.WithAuth(sqlc.RoleAdmin, context, createRecipe.Post(context)))

	mux.HandleFunc("GET /recipes/{id}", recipe.Get(context))

	mux.HandleFunc("GET /login", login.Get(context))
	mux.HandleFunc("POST /login", login.Post(context))

	mux.HandleFunc("GET /register", register.Get(context))
	mux.HandleFunc("POST /register", register.Post(context))

	mux.HandleFunc("GET /logout", auth.WithAuth(sqlc.RoleStandard, context, logout.Get(context)))
	mux.HandleFunc("POST /logout", auth.WithAuth(sqlc.RoleStandard, context, logout.Post(context)))

	mux.HandleFunc("GET /list", auth.WithAuth(sqlc.RoleStandard, context, list.Get(context)))
	mux.HandleFunc("POST /list", auth.WithAuth(sqlc.RoleStandard, context, list.Post(context)))

	mux.HandleFunc("GET /ingredients", ingredients.Get(context))
	mux.HandleFunc("POST /ingredients", ingredients.Post(context))

	mux.HandleFunc("GET /recipeCategories", auth.WithAuth(sqlc.RoleAdmin, context, recipeCategories.Get(context)))
	mux.HandleFunc("POST /recipeCategories", auth.WithAuth(sqlc.RoleAdmin, context, recipeCategories.Post(context)))

}
