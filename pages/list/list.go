package list

import (
	"fmt"
	"josephwest2/meal-list/assert"
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"josephwest2/meal-list/lib/db"
	"josephwest2/meal-list/pages"
	"net/http"
	"strconv"
)

func Get(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserIDFromSession(context.DB, r)
		assert.Assert(err == nil, "UserID is null in WithAuth protected path")
		list := context.DB.GetOrCreateListByUserID(userID)
		pages.RenderPage("List", List(&list), nil, w, r)
	}
}
func Post(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserIDFromSession(context.DB, r)
		assert.Assert(err == nil, "UserID is null in WithAuth protected path")
		list := context.DB.GetOrCreateListByUserID(userID)
		messages := make([]pages.PageMessage, 0)
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
		messages = append(messages, pages.PageMessage{Type: pages.Success, Value: "Item added to list", Timeout: true})
		list = context.DB.GetOrCreateListByUserID(userID)
		pages.RenderPage("List", List(&list), messages, w, r)
	}
}
