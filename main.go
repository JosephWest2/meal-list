package main

import (
	"josephwest2/meal-list/api/seed"
	"josephwest2/meal-list/lib/db"
	"josephwest2/meal-list/pages/index"
	"josephwest2/meal-list/pages/recipes"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}
	// needs to be called early
	db.Setup()

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// pages
	mux.HandleFunc("/", index.Handler)
	mux.HandleFunc("/recipes", recipes.Handler)

	// apis
	mux.HandleFunc("/api/seed", seed.Handler)

	http.ListenAndServe(":3000", mux)
}
