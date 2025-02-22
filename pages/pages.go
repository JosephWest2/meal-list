package pages

import (
	"context"
	"josephwest2/meal-list/lib/auth"
	"net/http"
	"strings"

	"github.com/a-h/templ"
)

type MessageType uint

const (
	Message MessageType = iota
	Success
	Warning
	Error
)

type PageMessage struct {
	Value   string
	Type    MessageType
	Timeout bool
}

func RedirectWithMessage(w http.ResponseWriter, r *http.Request, path string, message PageMessage) {
	s := strings.ReplaceAll(message.Value, " ", "+")
	switch message.Type {
	case Message:
		s = "?message=" + s
	case Success:
		s = "?success=" + s
	case Warning:
		s = "?warning=" + s
	case Error:
		s = "?error=" + s
	}
	http.Redirect(w, r, path+s, http.StatusSeeOther)
}

func RenderPage(pageTitle string, pageComponent templ.Component, messages []PageMessage, w http.ResponseWriter, r *http.Request) {

	messageQuery := r.URL.Query()["message"]
	if messageQuery != nil {
		messages = append(messages, PageMessage{Type: Message, Value: messageQuery[0]})
	}
	warningQuery := r.URL.Query()["warning"]
	if warningQuery != nil {
		messages = append(messages, PageMessage{Type: Warning, Value: messageQuery[0]})
	}
	successQuery := r.URL.Query()["success"]
	if successQuery != nil {
		messages = append(messages, PageMessage{Type: Success, Value: messageQuery[0]})
	}
	errorQuery := r.URL.Query()["error"]
	if errorQuery != nil {
		messages = append(messages, PageMessage{Type: Error, Value: messageQuery[0]})
	}
	isLoggedIn := auth.IsLoggedInUnverified(r)
	page := Layout(pageTitle, messages, isLoggedIn, pageComponent)
	page.Render(context.Background(), w)
}
