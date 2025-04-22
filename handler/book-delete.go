package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) deleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid URL", http.StatusInternalServerError)
		return
	}

	const getBook = `SELECT * FROM books WHERE id=$1`
	var book Books
	h.db.Get(&book, getBook, id)

	if book.ID == 0 {
		http.Error(w, "Invalid URL", http.StatusInternalServerError)
		return
	}
	const deleteBook = `DELETE FROM books WHERE id=$1`
	res := h.db.MustExec(deleteBook, id)
	ok, err := res.RowsAffected()
	if err != nil || ok == 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
