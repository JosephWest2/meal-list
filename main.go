package main

import (
	"josephwest2/meal-list/api/addRecipeToList"
	"josephwest2/meal-list/api/seed"
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"josephwest2/meal-list/lib/db"
	"josephwest2/meal-list/pages/index"
	"josephwest2/meal-list/pages/list"
	"josephwest2/meal-list/pages/login"
	"josephwest2/meal-list/pages/logout"
	"josephwest2/meal-list/pages/recipes"
	createRecipe "josephwest2/meal-list/pages/recipes/create"
	"josephwest2/meal-list/pages/recipes/recipe"
	"josephwest2/meal-list/pages/register"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env")
	}
	postgresConnectionString := os.Getenv("POSTGRES_CONNECTION_STRING")
	dbRef := db.InitDB(postgresConnectionString)

	context := app.AppContext{DB: &dbRef}

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// pages
	mux.HandleFunc("GET /", index.Handler)

	mux.HandleFunc("GET /recipes", recipes.Get(&context))
	mux.HandleFunc("POST /recipes", auth.WithAuth(auth.AdminRole, &context, recipes.Post(&context)))

	mux.HandleFunc("GET /recipes/create", auth.WithAuth(auth.AdminRole, &context, createRecipe.Get(&context)))
	mux.HandleFunc("POST /recipes/create", auth.WithAuth(auth.AdminRole, &context, createRecipe.Post(&context)))

	mux.HandleFunc("GET /recipes/{id}", recipe.Get(&context))

	mux.HandleFunc("GET /login", login.Get(&context))
	mux.HandleFunc("POST /login", login.Post(&context))

	mux.HandleFunc("GET /register", register.Get(&context))
	mux.HandleFunc("POST /register", register.Post(&context))

	mux.HandleFunc("GET /logout", auth.WithAuth(auth.StandardUserRole, &context, logout.Get(&context)))
	mux.HandleFunc("POST /logout", auth.WithAuth(auth.StandardUserRole, &context, logout.Post(&context)))

	mux.HandleFunc("GET /list", auth.WithAuth(auth.StandardUserRole, &context, list.Get(&context)))
	mux.HandleFunc("POST /list", auth.WithAuth(auth.StandardUserRole, &context, list.Post(&context)))

	// apis
	mux.HandleFunc("POST /api/seed", auth.WithAuth(auth.AdminRole, &context, seed.Post(&context)))
	mux.HandleFunc("POST /api/addRecipeToList/{id}", auth.WithAuth(auth.StandardUserRole, &context, addRecipeToList.Post(&context)))

	http.ListenAndServe(":3000", mux)
}
