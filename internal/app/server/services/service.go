package services

import (
	"context"

	"github.com/Te8va/APIbook/internal/app/server/domain"
)

type Book struct {
	repo domain.BookRepository
}

func NewBookService(repo domain.BookRepository) *Book {
	return &Book{repo: repo}
}

func (s *Book) GetBookByID(id string) (domain.Book, error) {
	return s.repo.GetBookByID(id)
}

func (s *Book) AddBook(ctx context.Context, newBook domain.Book) error {
	return s.repo.AddBook(ctx, newBook)
}

func (s *Book) DeleteBook(ctx context.Context, id string) error {
	return s.repo.DeleteBook(ctx, id)
}

func (s *Book) UpdateBook(ctx context.Context, id string, updatedBook domain.Book) error {
	return s.repo.UpdateBook(ctx, id, updatedBook)
}
