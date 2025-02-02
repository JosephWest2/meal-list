package index

import (
	"josephwest2/meal-list/pages"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	pages.RenderPage("Home", Index(r.URL.Path[1:]), nil, w, r)
}
