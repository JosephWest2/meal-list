package seed

import (
	"fmt"
	"josephwest2/meal-list/lib/db"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Seed request received")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	db.SeedRecipes()
	w.WriteHeader(http.StatusCreated)
}
