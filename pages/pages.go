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
	Value string
	Type  MessageType
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
	messageParam := ""
	if messageQuery != nil {
		messageParam = messageQuery[0]
	}
	warningQuery := r.URL.Query()["warning"]
	warningParam := ""
	if warningQuery != nil {
		warningParam = warningQuery[0]
	}
	successQuery := r.URL.Query()["success"]
	successParam := ""
	if successQuery != nil {
		successParam = successQuery[0]
	}
	errorQuery := r.URL.Query()["error"]
	errorParam := ""
	if errorQuery != nil {
		errorParam = errorQuery[0]
	}
	for _, message := range messages {
		switch message.Type {
		case Message:
			messageParam = message.Value
		case Success:
			successParam = message.Value
		case Warning:
			warningParam = message.Value
		case Error:
			errorParam = message.Value
		}
	}
	isLoggedIn := auth.IsLoggedInUnverified(r)
	page := Layout(pageTitle, messageParam, warningParam, successParam, errorParam, isLoggedIn, pageComponent)
	page.Render(context.Background(), w)
}
