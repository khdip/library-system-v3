package handler

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

const sessionName = "library-session"

type Books struct {
	ID          int    `db:"id" json:"id"`
	BookName    string `db:"book_name" json:"book_name"`
	Author      string `db:"author" json:"author"`
	Category    string `db:"category" json:"category"`
	BookDesc    string `db:"book_description" json:"book_description"`
	BookCover   string `db:"book_cover" json:"book_cover"`
	IsAvailable bool   `db:"is_available" json:"is_available"`
}

type Handler struct {
	templates *template.Template
	db        *sqlx.DB
	decoder   *schema.Decoder
	session   *sessions.CookieStore
}

func GetHandler(db *sqlx.DB, decoder *schema.Decoder, session *sessions.CookieStore) *mux.Router {
	hand := &Handler{
		db:      db,
		decoder: decoder,
		session: session,
	}
	hand.GetTemplate()

	r := mux.NewRouter()
	r.HandleFunc("/", hand.GetBooks)

	loginRouter := r.NewRoute().Subrouter()
	loginRouter.HandleFunc("/login", hand.login)
	loginRouter.HandleFunc("/login/auth", hand.loginAuth)
	loginRouter.HandleFunc("/signup", hand.signup)
	loginRouter.HandleFunc("/register", hand.register)
	loginRouter.Use(hand.restrictMiddleware)
	r.HandleFunc("/logout", hand.logout)

	s := r.NewRoute().Subrouter()
	s.HandleFunc("/create", hand.createBook)
	s.HandleFunc("/store", hand.storeBook)
	r.HandleFunc("/q", hand.searchBook)
	s.HandleFunc("/{id:[0-9]+}/edit", hand.editBook)
	s.HandleFunc("/{id:[0-9]+}/Update", hand.updateBook)
	s.HandleFunc("/{id:[0-9]+}/delete", hand.deleteBook)
	s.Use(hand.authMiddleware)
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := hand.templates.ExecuteTemplate(w, "404.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	return r
}

func (h *Handler) GetTemplate() {
	h.templates = template.Must(template.ParseFiles(
		"templates/create-book.html",
		"templates/list-book.html",
		"templates/edit-book.html",
		"templates/search-result.html",
		"templates/no-search-result.html",
		"templates/404.html",
		"templates/login.html",
		"templates/signup.html",
	))
}
