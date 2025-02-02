package register

import (
	"errors"
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/db"
	"josephwest2/meal-list/pages"
	"net/http"
	"strings"
)

func IsValidPassword(password string) []error {
	errs := make([]error, 0)
	if len(password) < 8 {
		errs = append(errs, errors.New("password must be at least 8 characters"))
	}
	if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		errs = append(errs, errors.New("password must contain at least one letter"))
	}
	if !strings.ContainsAny(password, "0123456789") {
		errs = append(errs, errors.New("password must contain at least one number"))
	}
	if !strings.ContainsAny(password, " !\"#$%&'()*+,-./:;<=>?@[]^_`{|}~\\") {
		errs = append(errs, errors.New("password must contain at least one special character"))
	}
	return errs
}

func Handler(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		messages := make([]pages.PageMessage, 0)
		switch r.Method {
		case "GET":
		case "POST":
			r.ParseForm()
			username := r.FormValue("username")
			password := r.FormValue("password")
			registrationSuccess := true
			if len(username) < 3 {
				messages = append(messages, pages.PageMessage{Type: pages.Error, Value: "Username must be at least 3 characters"})
				registrationSuccess = false
				w.WriteHeader(http.StatusBadRequest)
			}
			_, err := context.DB.GetUserByUsername(username)
			if err == nil {
				messages = append(messages, pages.PageMessage{Type: pages.Error, Value: "Username taken"})
				w.WriteHeader(http.StatusConflict)
				registrationSuccess = false
			}
			errs := IsValidPassword(password)
			if len(errs) > 0 {
				registrationSuccess = false
				w.WriteHeader(http.StatusBadRequest)
				for _, err := range errs {
					messages = append(messages, pages.PageMessage{Type: pages.Error, Value: err.Error()})
				}
			}
			if !registrationSuccess {
				break
			}
			err = context.DB.CreateUser(username, password, db.StandardUser)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			pages.RedirectWithMessage(w, r, "/login", pages.PageMessage{Type: pages.Success, Value: "Registration Success"})
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		pages.RenderPage("Register", Register(), messages, w, r)
	}
}
