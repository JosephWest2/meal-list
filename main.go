package main

import (
	"josephwest2/meal-list/api/seed"
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"josephwest2/meal-list/lib/db"
	"josephwest2/meal-list/pages/index"
	"josephwest2/meal-list/pages/list"
	"josephwest2/meal-list/pages/login"
	"josephwest2/meal-list/pages/logout"
	"josephwest2/meal-list/pages/recipes"
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
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// pages
	mux.HandleFunc("/", index.Handler)
	mux.HandleFunc("/recipes", recipes.Handler(&context))
	mux.HandleFunc("/recipes/", recipe.Handler(&context))
	mux.HandleFunc("/login", login.Handler(&context))
	mux.HandleFunc("/register", register.Handler(&context))
	mux.HandleFunc("/logout", auth.WithAuth(db.StandardUser, &context, logout.Handler(&context)))
	mux.HandleFunc("/list", auth.WithAuth(db.StandardUser, &context, list.Handler(&context)))

	// apis
	mux.HandleFunc("/api/seed", seed.Handler(&context))

	http.ListenAndServe(":3000", mux)
}
