package seed

import (
	"fmt"
	"josephwest2/meal-list/lib/app"
	"net/http"
)

func Handler(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Seed request received")
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		context.DB.Seed()
		w.WriteHeader(http.StatusCreated)

	}
}
