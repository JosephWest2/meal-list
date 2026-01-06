package main

import (
	"context"
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/router"
	"josephwest2/meal-list/lib/sqlc"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env")
	}
	postgresConnectionString := os.Getenv("POSTGRES_CONNECTION_STRING")
	pool, err := pgxpool.New(ctx, postgresConnectionString)
	queries := sqlc.New(pool)

	context := app.AppContext{Queries: queries, QueryContext: ctx}

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	router.RegisterAPIRoutes(mux, &context)
	router.RegisterPageRoutes(mux, &context)

	http.ListenAndServe(":3000", mux)
}
