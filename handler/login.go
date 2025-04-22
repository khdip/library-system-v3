package handler

import (
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

type LoginFormData struct {
	Email    string
	Password string
	Errors   map[string]string
	Notices  string
}

func (l LoginFormData) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Email, validation.Required.Error("The email field is required")),
		validation.Field(&l.Password, validation.Required.Error("The password field is required"), validation.Length(6, 20).Error("Password must be between 6 to 20 characters")),
	)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	session, err := h.session.Get(r, sessionName)
	if err != nil {
		log.Fatal(err)
	}

	form := LoginFormData{}

	if flashes := session.Flashes(); len(flashes) > 0 {
		if val, ok := flashes[0].(string); ok {
			form.Notices = val
		}
	}

	if err := session.Save(r, w); err != nil {
		log.Fatal(err)
	}

	h.loadLoginForm(w, form)
}

func (h *Handler) loginAuth(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var form LoginFormData
	err = h.decoder.Decode(&form, r.PostForm)
	if err != nil {
		log.Fatal(err)
	}

	err = form.Validate()
	if err != nil {
		vErrors, ok := err.(validation.Errors)
		if ok {
			ErrorValue := make(map[string]string)
			for _, value := range vErrors {
				ErrorValue = map[string]string{
					"Error": value.Error(),
				}
			}
			form.Errors = ErrorValue
			h.loadLoginForm(w, form)
			return
		}
	}

	const getUser = "SELECT * FROM users WHERE email=$1"
	var user UserDetails
	h.db.Get(&user, getUser, form.Email)
	if user.UserID == 0 {
		form.Errors = map[string]string{"Error": "Invalid User"}
		h.loadLoginForm(w, form)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		form.Errors = map[string]string{"Error": "Wrong Password"}
		h.loadLoginForm(w, form)
		return
	}

	session, err := h.session.Get(r, sessionName)
	if err != nil {
		log.Fatal(err)
	}

	session.Options.HttpOnly = true

	session.Values["authUserId"] = user.UserID
	if err := session.Save(r, w); err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (h *Handler) loadLoginForm(w http.ResponseWriter, form LoginFormData) {
	err := h.templates.ExecuteTemplate(w, "login.html", form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
