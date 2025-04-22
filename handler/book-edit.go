package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

func (h *Handler) editBook(w http.ResponseWriter, r *http.Request) {
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
	h.loadEditForm(w, book, map[string]string{})
}

func (h *Handler) updateBook(w http.ResponseWriter, r *http.Request) {
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

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.decoder.Decode(&book, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = book.Validate()
	if err != nil {
		vErrors, ok := err.(validation.Errors)
		if ok {
			ErrorValue := make(map[string]string)
			for _, value := range vErrors {
				ErrorValue = map[string]string{
					"Error": value.Error(),
				}
			}
			h.loadCreateForm(w, book, ErrorValue)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("BookCover")
	if err != nil {
		ErrorValue := map[string]string{
			"Error": "Book Cover is required & image size should be less than 10MB.",
		}
		h.loadEditForm(w, book, ErrorValue)
		return
	}
	defer file.Close()
	tempFile, err := ioutil.TempFile("assets/book-covers", "cover-*.jpg")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)
	bookCoverPath := tempFile.Name()
	book.BookCover = strings.TrimPrefix(bookCoverPath, "assets/book-covers\\")

	const updateTodo = `UPDATE books SET book_name = $2, author = $3, category = $4, book_description = $5, book_cover = $6, is_available = $7 WHERE id=$1`
	res := h.db.MustExec(updateTodo, id, book.BookName, book.Author, book.Category, book.BookDesc, book.BookCover, book.IsAvailable)
	ok, err := res.RowsAffected()
	if err != nil || ok == 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
