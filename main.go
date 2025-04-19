package main

import (
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/db"
	"josephwest2/meal-list/lib/router"
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

	router.RegisterAPIRoutes(mux, &context)
	router.RegisterPageRoutes(mux, &context)


	http.ListenAndServe(":3000", mux)
}
