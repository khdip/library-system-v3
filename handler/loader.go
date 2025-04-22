package handler

import "net/http"

type FormData struct {
	Book   Books
	Errors map[string]string
}

func (h *Handler) loadCreateForm(w http.ResponseWriter, book Books, myErrors map[string]string) {
	form := FormData{
		Book:   book,
		Errors: myErrors,
	}

	err := h.templates.ExecuteTemplate(w, "create-book.html", form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) loadEditForm(w http.ResponseWriter, book Books, myErrors map[string]string) {
	form := FormData{
		Book:   book,
		Errors: myErrors,
	}

	err := h.templates.ExecuteTemplate(w, "edit-book.html", form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
