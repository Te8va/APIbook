package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"github.com/Te8va/APIbook/internal/app/server/domain"
	logging "github.com/Te8va/APIbook/internal/pkg/logger"
)

type Book struct {
	srv domain.BookService
}

func NewBookHandler(srv domain.BookRepository) *Book {
	return &Book{srv: srv}
}

func (h *Book) GetBookByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	if uuid.Validate(id) != nil {
		logging.Logger().Warn("Invalid book ID format")
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	book, err := h.srv.GetBookByID(id)
	reply(w, book, http.StatusOK, err)
}

func (h *Book) AddBookHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var newBook domain.Book
	if err := d.Decode(&newBook); err != nil {
		logging.Logger().Error("Error decoding request payload", err)
		http.Error(w, "Invalid to decode request payload", http.StatusBadRequest)
		return
	}

	if newBook.Title == "" || newBook.Author == "" || newBook.Year == 0 {
		logging.Logger().Warn("Required fields are not filled in for adding a new book: Title, Author, Year")
		http.Error(w, "Title, Author, and Data are required fields", http.StatusBadRequest)
		return
	}

	bookID, err := h.srv.AddBook(r.Context(), newBook)
	reply(w, bookID, http.StatusCreated, err)
}

func (h *Book) DeleteBookHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	if uuid.Validate(id) != nil {
		logging.Logger().Warn("Invalid book ID format")
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err := h.srv.DeleteBook(r.Context(), id)
	reply(w, "Book deleted successfully", http.StatusOK, err)
}

func (h *Book) UpdateBookHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var updatedBook domain.Book
	err := d.Decode(&updatedBook)
	if err != nil {
		logging.Logger().Error("Error decoding request payload", err)
		http.Error(w, "Invalid to decode request payload", http.StatusBadRequest)
		return
	}

	updatedBook.ID = ps.ByName("id")

	if uuid.Validate(updatedBook.ID) != nil {
		logging.Logger().Warn("Invalid book ID format")
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if updatedBook.Title == "" || updatedBook.Author == "" || updatedBook.Year == 0 {
		logging.Logger().Warn("Required fields are not filled in for adding a new book: Title, Author, Year")
		http.Error(w, "Title, Author, and Data are required fields", http.StatusBadRequest)
		return
	}

	err = h.srv.UpdateBook(r.Context(), updatedBook.ID, updatedBook)
	reply(w, "Updated", http.StatusOK, err)
}
