package handler

import (
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type SignupFormData struct {
	UserDetails UserDetails
	Errors      map[string]string
}

type UserDetails struct {
	UserID          int    `db:"user_id"`
	FirstName       string `db:"first_name"`
	LastName        string `db:"last_name"`
	Email           string `db:"email"`
	Password        string `db:"password"`
	ConfirmPassword string
	IsVerified      bool `db:"is_verified"`
}

func (u UserDetails) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.FirstName, validation.Required.Error("The first name field is required")),
		validation.Field(&u.LastName, validation.Required.Error("The last name field is required")),
		validation.Field(&u.Email, validation.Required.Error("The email field is required"), is.Email),
		validation.Field(&u.Password, validation.Required.Error("The password field is required"), validation.Length(6, 20).Error("Password must be between 6 to 20 characters")),
		validation.Field(&u.ConfirmPassword, validation.Required.Error("The confirm password field is required"), validation.Length(6, 20).Error("Password must be between 6 to 20 characters")),
	)
}

func (h *Handler) signup(w http.ResponseWriter, r *http.Request) {
	formData := SignupFormData{}
	h.loadSignupForm(w, formData)
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var form UserDetails
	err = h.decoder.Decode(&form, r.PostForm)
	if err != nil {
		log.Fatal(err)
	}

	if form.Password != form.ConfirmPassword {
		formData := SignupFormData{
			UserDetails: form,
			Errors:      map[string]string{"Error": "Passwords didn't match"},
		}
		h.loadSignupForm(w, formData)
		return
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
			formData := SignupFormData{
				UserDetails: form,
				Errors:      ErrorValue,
			}
			h.loadSignupForm(w, formData)
			return
		}
	}

	const insertUser = `INSERT INTO users(first_name, last_name, email, password, is_verified) VALUES($1, $2, $3, $4, $5);`
	passByte, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	res := h.db.MustExec(insertUser, form.FirstName, form.LastName, form.Email, string(passByte), false)
	ok, err := res.RowsAffected()
	if err != nil || ok == 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, err := h.session.Get(r, sessionName)
	if err != nil {
		log.Fatal(err)
	}
	session.AddFlash("Registration Successfull")
	if err := session.Save(r, w); err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

func (h *Handler) loadSignupForm(w http.ResponseWriter, formData SignupFormData) {
	err := h.templates.ExecuteTemplate(w, "signup.html", formData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
