package handler

import validation "github.com/go-ozzo/ozzo-validation"

func (book *Books) Validate() error {
	return validation.ValidateStruct(book,
		validation.Field(&book.BookName, validation.Required.Error("Book Name field can not be empty."), validation.Length(3, 50).Error("Book Name field should have atleast 3 characters and atmost 50 characters")),
	)
}
