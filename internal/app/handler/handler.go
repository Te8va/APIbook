package handler

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/Te8va/APIbook/internal/app/domain"
)

type Book struct {
	srv domain.BookService
}

func NewBookHandler(srv domain.BookRepository) *Book {
	return &Book{srv: srv}
}

func (h *Book) GetBookByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: use r.Context if you don't need to swap it with some other context/want to cancel e.t.c.
	ctx := r.Context()

	id := ps.ByName("id")

	book, err := h.srv.GetBookByID(ctx, id)
	if err != nil {
		handleBookError(w, err)
		return
	}

	// TODO: remove this, the case should be handled by handleError
	if book.ID == "" {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	// TODO: remove redundant comments
	// if err := reply(w, book); err != nil {
	// 	http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	// }

	reply(w, book, http.StatusOK)
}

func (h *Book) AddBookHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	defer r.Body.Close()
	// TODO: use disallow unknown fields

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var newBook domain.Book
	if err := d.Decode(&newBook); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest) // TODO: change error description
		return
	}

	if newBook.Title == "" || newBook.Author == "" || newBook.Year == 0 {
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

	if err := h.srv.DeleteBook(ctx, id); err != nil {
		handleBookError(w, err)
		return
	}

	reply(w, "Created", http.StatusOK) // TODO: check if status OK is default value, change message to actual one
}

func (h *Book) UpdateBookHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	id := ps.ByName("id") // TODO: you can get the id from book in body

	defer r.Body.Close()

	var updatedBook domain.Book
	err := json.NewDecoder(r.Body).Decode(&updatedBook)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if updatedBook.Title == "" || updatedBook.Author == "" || updatedBook.Year == 0 {
		http.Error(w, "Title, Author, and Data are required fields", http.StatusBadRequest)
		return
	}

	if err = h.srv.UpdateBook(ctx, id, updatedBook); err != nil { // TODO: use updatedBook.ID
		handleBookError(w, err)
		return
	}

	reply(w, "Created", http.StatusOK)
}
