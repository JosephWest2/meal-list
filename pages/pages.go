package pages

import (
	"context"
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"josephwest2/meal-list/lib/sqlc"
	"net/http"
	"strings"

	"github.com/a-h/templ"
)

type MessageType uint

const (
	MessageInfo MessageType = iota
	MessageSuccess
	MessageWarning
	MessageError
)

type PageMessage struct {
	Value   string
	Type    MessageType
	Timeout bool
}

func RedirectWithMessage(w http.ResponseWriter, r *http.Request, path string, message PageMessage) {
	s := strings.ReplaceAll(message.Value, " ", "+")
	switch message.Type {
	case MessageInfo:
		s = "?message=" + s
	case MessageSuccess:
		s = "?success=" + s
	case MessageWarning:
		s = "?warning=" + s
	case MessageError:
		s = "?error=" + s
	}
	http.Redirect(w, r, path+s, http.StatusSeeOther)
}

func RenderPage(app *app.App, pageTitle string, pageComponent templ.Component, messages []PageMessage, w http.ResponseWriter, r *http.Request) {

	messageQuery := r.URL.Query()["message"]
	if messageQuery != nil {
		messages = append(messages, PageMessage{Type: MessageInfo, Value: messageQuery[0]})
	}
	warningQuery := r.URL.Query()["warning"]
	if warningQuery != nil {
		messages = append(messages, PageMessage{Type: MessageWarning, Value: warningQuery[0]})
	}
	successQuery := r.URL.Query()["success"]
	if successQuery != nil {
		messages = append(messages, PageMessage{Type: MessageSuccess, Value: successQuery[0]})
	}
	errorQuery := r.URL.Query()["error"]
	if errorQuery != nil {
		messages = append(messages, PageMessage{Type: MessageError, Value: errorQuery[0]})
	}
	isLoggedIn := auth.IsAuthenticated(app, r)
	isAdmin := auth.IsAuthorized(app, r, sqlc.RoleAdmin)
	page := Layout(pageTitle, messages, isLoggedIn, isAdmin, pageComponent)
	page.Render(context.Background(), w)
}
