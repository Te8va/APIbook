package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"

	"github.com/Te8va/APIbook/internal/app/server/domain"
)

var _ domain.BookRepository = (*Book)(nil)

type Book struct {
	FilePath string
	mu       *sync.RWMutex
}

func NewFileBookRepository(filePath string) (*Book, error) {
	err := os.MkdirAll(filePath, 0755)
	if err != nil {
		return nil, domain.ErrCreatingDirectory
	}

	return &Book{FilePath: filePath}, nil
}

func (f *Book) GetBookByID(id string) (domain.Book, error) {
	filePath := fmt.Sprintf("%s/book_%s.json", f.FilePath, id)

	f.mu.RLock()
	defer f.mu.RUnlock()

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return domain.Book{}, domain.ErrBookNotFound
		}

		return domain.Book{}, domain.ErrReadingFile
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var book domain.Book
	err = decoder.Decode(&book)
	if err != nil {
		return domain.Book{}, domain.ErrDecodingJSON
	}

	if book.IsDeleted {
		return domain.Book{}, domain.ErrDeletedBook
	}

	return book, nil
}

func (f *Book) AddBook(ctx context.Context, newBook domain.Book) (string, error) {
	newUUID := uuid.New()

	fileName := fmt.Sprintf("book_%s.json", newUUID)

	f.mu.RLock()
	defer f.mu.RUnlock()

	filePath := filepath.Join(f.FilePath, fileName)

	newBook.ID = newUUID.String()

	data, err := json.MarshalIndent(newBook, "", "    ")
	if err != nil {
		return "", domain.ErrEncodingJSON
	}

	err = os.WriteFile(filePath, data, 0660)
	if err != nil {
		return "", domain.ErrWritingToFile
	}

	return newUUID.String(), nil
}

func (f *Book) UpdateBook(ctx context.Context, id string, updatedBook domain.Book) error {
	filePath := fmt.Sprintf("%s/book_%s.json", f.FilePath, id)

	f.mu.RLock()
	defer f.mu.RUnlock()

	file, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return domain.ErrBookNotFound
		}

		return domain.ErrReadingFile
	}

	var book domain.Book
	if err := json.Unmarshal(file, &book); err != nil {
		return domain.ErrDecodingJSON
	}

	if book.IsDeleted {
		return domain.ErrDeletedBook
	}

	updatedBook.ID = book.ID

	data, err := json.MarshalIndent(updatedBook, "", "    ")
	if err != nil {
		return domain.ErrEncodingJSON
	}

	err = os.WriteFile(filePath, data, 0660)
	if err != nil {
		return domain.ErrWritingToFile
	}

	return nil
}

func (f *Book) DeleteBook(ctx context.Context, id string) error {
	filePath := fmt.Sprintf("%s/book_%s.json", f.FilePath, id)

	f.mu.Lock()
	defer f.mu.Unlock()

	file, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return domain.ErrBookNotFound
		}

		return domain.ErrReadingFile
	}

	var book domain.Book
	if err := json.Unmarshal(file, &book); err != nil {
		return domain.ErrDecodingJSON
	}

	if book.IsDeleted {
		return domain.ErrDeletedBook
	}

	book.IsDeleted = true

	data, err := json.MarshalIndent(book, "", "    ")
	if err != nil {
		return domain.ErrEncodingJSON
	}

	err = os.WriteFile(filePath, data, 0660)
	if err != nil {
		return domain.ErrWritingToFile
	}

	return nil
}
