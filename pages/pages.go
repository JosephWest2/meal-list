package pages

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
)

func RenderPage(pageTitle string, pageComponent templ.Component, w http.ResponseWriter, r *http.Request) {
	messageParam := r.URL.Query()["message"]
	message := ""
	if messageParam != nil {
		message = messageParam[0]
	}
	page := Layout(pageTitle, message, pageComponent)
	page.Render(context.Background(), w)
}
