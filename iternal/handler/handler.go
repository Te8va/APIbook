package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"fmt"
	"io/fs"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/Te8va/APIbook.git/iternal/domain"
)

type BookHandler struct {
	srv domain.Book
}

var books []domain.Book

func NewBookHandler(srv domain.Book) *BookHandler {
	return &BookHandler{srv: srv}
}


const filePath = "books.json" //TODO файловую часть в другой пакет?

func ReadBooksFromFile() error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("Error reading file: %v", err)
	}

	err = json.Unmarshal(file, &books)
	if err != nil {
		return fmt.Errorf("Error decoding JSON: %v", err)
	}
	return nil
}

func SaveBooksToFile() error{
	data, err := json.Marshal(books)
	if err != nil {
		return fmt.Errorf("Error encoding JSON: %v", err)
	}

	err = os.WriteFile(filePath, data, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("Error writing to file:%v", err)
	}
	return nil
}

func (bh *BookHandler) GetBooksHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    if err := reply(w, books); err != nil {
        http.Error(w, "Error retrieving books", http.StatusInternalServerError)
    }
}

func (bh *BookHandler) GetBookByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	for _, book := range books {
		if book.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(book)
			return
		}
	}

	http.Error(w, "Book not found", http.StatusNotFound)
}

func (bh *BookHandler) AddBookHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var newBook domain.Book

	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	newBook.ID = len(books) + 1

	books = append(books, newBook)

	
	if err := SaveBooksToFile(); err != nil {
		http.Error(w, "Error saving to file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newBook)
}

func (bh *BookHandler) DeleteBookHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := validateID(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	index := -1
	for i, book := range books {
		if book.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	books = append(books[:index], books[index+1:]...)

	if err := SaveBooksToFile(); err != nil {
		http.Error(w, "Error saving to file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (bh *BookHandler) UpdateBookHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := validateID(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	var updatedBook domain.Book
	err = json.NewDecoder(r.Body).Decode(&updatedBook)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedBook.ID = id

	index := -1
	for i, book := range books {
		if book.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	books[index] = updatedBook

	err = SaveBooksToFile()
	if err != nil {
		http.Error(w, "Error saving to file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func validateID(param string) (int, error) {
	id, err := strconv.Atoi(param)
	if err != nil {
		return 0, fmt.Errorf("Invalid ID: %v", err)
	}
	return id, nil
}

func reply(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
        return fmt.Errorf("Error encoding JSON: %v", err)
    }
	return nil
}