package list

import (
	"fmt"
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"josephwest2/meal-list/lib/db"
	"josephwest2/meal-list/pages"
	"net/http"
	"strconv"
)

func Handler(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserIDFromSession(context.DB, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		list := context.DB.GetListByUserID(userID)
		switch r.Method {
		case "GET":
		case "POST":
			r.ParseForm()
			fmt.Println(r)
			name := r.Form.Get("name")
			quanityRaw := r.Form.Get("quantity")
			quantity, err := strconv.ParseFloat(quanityRaw, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			unit := r.Form.Get("unit")
			listItemDetails := db.ListItemDetails{
				Name:     name,
				Quantity: quantity,
				Unit:     unit,
			}
			err = context.DB.AddToList(list.ID, listItemDetails)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		case "PATCH":
		case "DELETE":
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		pages.RenderPage("List", List(&list), nil, w, r)
	}
}
