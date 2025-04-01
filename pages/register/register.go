package register

import (
	"errors"
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
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

func Get(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pages.RenderPage("Register", Register(), nil, w, r)
	}
}

func Post(context *app.AppContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		messages := make([]pages.PageMessage, 0)
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
            pages.RenderPage("Register", Register(), messages, w, r)
            return
        }
        err = context.DB.CreateUser(username, password, auth.StandardUserRole)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        pages.RedirectWithMessage(w, r, "/login", pages.PageMessage{Type: pages.Success, Value: "Registration Success"})
        return

    }
}
