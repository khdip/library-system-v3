package handler

import (
	"net/http"
)

type BookList struct {
	Book_list []Books
}

func (h *Handler) GetBooks(w http.ResponseWriter, r *http.Request) {
	books := []Books{}
	h.db.Select(&books, "SELECT * FROM books")
	lt := BookList{Book_list: books}
	err := h.templates.ExecuteTemplate(w, "list-book.html", lt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
