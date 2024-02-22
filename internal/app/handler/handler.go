package handler

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"github.com/Te8va/APIbook/internal/app/domain"
	"github.com/Te8va/APIbook/internal/app/logging"
)

type Book struct {
	srv domain.BookService
}

func NewBookHandler(srv domain.BookRepository) *Book {
	return &Book{srv: srv}
}

func (h *Book) GetBookByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id := ps.ByName("id")

	matched, err := regexp.MatchString(`^[a-zA-Z0-9]+$`, id)
	if err != nil {
		logging.Logger.Error("Error validating book ID", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !matched {
		logging.Logger.Warn("Invalid book ID format")
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	book, err := h.srv.GetBookByID(id)
	if err != nil {
		handleBookError(w, err)
		return
	}

	reply(w, book, http.StatusOK)
}

func (h *Book) AddBookHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	defer r.Body.Close()

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var newBook domain.Book
	if err := d.Decode(&newBook); err != nil {
		logging.Logger.Error("Error decoding request payload", zap.Error(err))
		http.Error(w, "Invalid to decode request payload", http.StatusBadRequest)
		return
	}

	if newBook.Title == "" || newBook.Author == "" || newBook.Year == 0 {
		logging.Logger.Warn("Required fields are not filled in for adding a new book: Title, Author, Year")
		http.Error(w, "Title, Author, and Data are required fields", http.StatusBadRequest)
		return
	}

	if err := h.srv.AddBook(ctx, newBook); err != nil {
		handleBookError(w, err)
		return
	}

	reply(w, "Created", http.StatusCreated)
}

func (h *Book) DeleteBookHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	id := ps.ByName("id")

	matched, err := regexp.MatchString(`^[a-zA-Z0-9]+$`, id)
	if err != nil {
		logging.Logger.Error("Error validating book ID", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !matched {
		logging.Logger.Warn("Invalid book ID format")
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if err := h.srv.DeleteBook(ctx, id); err != nil {
		handleBookError(w, err)
		return
	}

	reply(w, "Book deleted successfully", http.StatusOK)
}

func (h *Book) UpdateBookHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	defer r.Body.Close()

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var updatedBook domain.Book
	err := d.Decode(&updatedBook)
	if err != nil {
		logging.Logger.Error("Error decoding request payload", zap.Error(err))
		http.Error(w, "Invalid to decode request payload", http.StatusBadRequest)
		return
	}

	updatedBook.ID = ps.ByName("id")

	matched, err := regexp.MatchString(`^[a-zA-Z0-9]+$`, updatedBook.ID)
	if err != nil {
		logging.Logger.Error("Error validating book ID", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !matched {
		logging.Logger.Warn("Invalid book ID format")
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if updatedBook.Title == "" || updatedBook.Author == "" || updatedBook.Year == 0 {
		logging.Logger.Warn("Required fields are not filled in for adding a new book: Title, Author, Year")
		http.Error(w, "Title, Author, and Data are required fields", http.StatusBadRequest)
		return
	}

	if err = h.srv.UpdateBook(ctx, updatedBook.ID, updatedBook); err != nil {
		handleBookError(w, err)
		return
	}

	reply(w, "Updated", http.StatusOK)
}
