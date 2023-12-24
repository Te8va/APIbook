package handler

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/Te8va/APIbook/docs"
	"github.com/Te8va/APIbook/internal/app/domain"
)

type Book struct {
	srv domain.BookService
}

func NewBookHandler(srv domain.BookRepository) *Book {
	return &Book{srv: srv}
}

func (h *Book) SwaggerHandler(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	httpSwagger.WrapHandler(res, req)
}

// @Tags Books
// @Summary Запрос чтения информации о книге
// @Description Запрос для получения сохраненной информации о книге
// @Produce json
// @Param id path string true "book id" Example(8502ab55-6750-4c53-8126-acc1ba19f801)
// @Success 200
// @Failure 404
// @Failure 500
// @Router /books/{id} [get]
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

// @Tags Books
// @Summary Запрос добавления книги
// @Description Запрос для добавления информации о новой книге
// @Accept json
// @Produce json
// @Param input body domain.Book true "book info"
// @Success 201
// @Failure 404
// @Failure 500
// @Failure 400
// @Router /books [post]
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
	id, err := h.srv.AddBook(ctx, newBook)
	if err != nil {
		handleBookError(w, err)
		return
	}

	reply(w, "Created, id: "+id, http.StatusCreated)
}

// @Tags Books
// @Summary Запрос удаления книги
// @Description Запрос для удаления информации о существующей книге
// @Produce json
// @Param id path string true "book id" Example(8502ab55-6750-4c53-8126-acc1ba19f801)
// @Success 200
// @Failure 404
// @Failure 500
// @Router /books/{id} [delete]
func (h *Book) DeleteBookHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	id := ps.ByName("id")

	if err := h.srv.DeleteBook(ctx, id); err != nil {
		handleBookError(w, err)
		return
	}

	reply(w, "Deleted", http.StatusOK) // TODO: check if status OK is default value, change message to actual one
}

// @Tags Books
// @Summary Запрос обновления информации о книге
// @Description Запрос для обновления информации о существующей книге
// @Accept json
// @Param input body domain.Book true "book info"
// @Param id path string true "book id" Example(8502ab55-6750-4c53-8126-acc1ba19f801)
// @Success 200
// @Failure 404
// @Failure 500
// @Failure 400
// @Router /books/{id} [put]
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

	reply(w, "Updated", http.StatusOK)
}
