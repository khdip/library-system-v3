package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (h *Handler) createBook(w http.ResponseWriter, r *http.Request) {
	ErrorValue := map[string]string{}
	book := Books{}
	h.loadCreateForm(w, book, ErrorValue)
}

func (h *Handler) storeBook(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var book Books
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
		h.loadCreateForm(w, book, ErrorValue)
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

	const insertBook = `INSERT INTO books(book_name, author, category, book_description, book_cover, is_available) VALUES($1, $2, $3, $4, $5, $6);`
	res := h.db.MustExec(insertBook, book.BookName, book.Author, book.Category, book.BookDesc, book.BookCover, book.IsAvailable)
	ok, err := res.RowsAffected()
	if err != nil || ok == 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
