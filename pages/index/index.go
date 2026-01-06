package index

import (
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/pages"
	"net/http"
)

func Get(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pages.RenderPage(context, "Home", Index(r.URL.Path[1:]), nil, w, r)
	}
}
