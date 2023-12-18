package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"

	"github.com/Te8va/APIbook/internal/app/domain"
)

type Book struct {
	FilePath string
	mu       *sync.RWMutex
}

func NewFileBookRepository(filePath string) *Book {
	return &Book{FilePath: filePath, mu: &sync.RWMutex{}}
}

func (r *Book) GetBookByID(ctx context.Context, id string) (domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filePath := fmt.Sprintf("%s/book_%s.json", r.FilePath, id) // TODO: check id in handler, it should only contain letters, numbers and -
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return domain.Book{}, domain.ErrBookNotFound
		}
		return domain.Book{}, fmt.Errorf(domain.ErrReadingFile.Error(), err) // TODO: read about error wrapping
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var book domain.Book
	err = decoder.Decode(&book)
	if err != nil {
		return domain.Book{}, fmt.Errorf(domain.ErrDecodingJSON.Error(), err)
	}

	if book.Status == "deleted" {
		return domain.Book{}, domain.ErrDeletedBook
	}

	return book, nil
}

func (f *Book) AddBook(ctx context.Context, updatedBook domain.Book) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	newUUID := uuid.New()

	fileName := fmt.Sprintf("book_%s.json", newUUID)
	filePath := filepath.Join(f.FilePath, fileName)

	updatedBook.ID = newUUID.String()

	data, err := json.MarshalIndent(updatedBook, "", "    ")
	if err != nil {
		return fmt.Errorf(domain.ErrEncodingJSON.Error(), err)
	}

	err = os.WriteFile(filePath, data, 0660)
	if err != nil {
		return fmt.Errorf(domain.ErrWritingToFile.Error(), err)
	}

	return nil
}

func (f *Book) UpdateBook(ctx context.Context, id string, updatedBook domain.Book) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	filePath := fmt.Sprintf("%s/book_%s.json", f.FilePath, id)
	if _, err := os.Stat(filePath); os.IsNotExist(err) { // TODO: check if it can be done with checking error of os.ReadFile
		return fmt.Errorf(domain.ErrBookNotFound.Error())
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf(domain.ErrReadingFile.Error(), err)
	}

	var books domain.Book
	if err := json.Unmarshal(file, &books); err != nil {
		return fmt.Errorf(domain.ErrDecodingJSON.Error(), err)
	}

	if books.Status == "deleted" {
		return domain.ErrDeletedBook
	}

	updatedBook.ID = books.ID

	data, err := json.MarshalIndent(updatedBook, "", "    ")
	if err != nil {
		return fmt.Errorf(domain.ErrEncodingJSON.Error(), err)
	}

	err = os.WriteFile(filePath, data, 0660)
	if err != nil {
		return fmt.Errorf(domain.ErrWritingToFile.Error(), err)
	}

	return nil
}

func (f *Book) DeleteBook(ctx context.Context, id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	filePath := fmt.Sprintf("%s/book_%s.json", f.FilePath, id)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return domain.ErrBookNotFound
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf(domain.ErrReadingFile.Error(), err)
	}

	var book domain.Book
	if err := json.Unmarshal(file, &book); err != nil {
		return fmt.Errorf(domain.ErrDecodingJSON.Error(), err)
	}

	if book.Status == "deleted" {
		return domain.ErrDeletedBook
	}

	book.Status = "deleted"

	data, err := json.MarshalIndent(book, "", "    ")
	if err != nil {
		return fmt.Errorf(domain.ErrEncodingJSON.Error(), err)
	}

	err = os.WriteFile(filePath, data, 0660)
	if err != nil {
		return fmt.Errorf(domain.ErrWritingToFile.Error(), err)
	}

	return nil
}
