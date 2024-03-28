package domain

import (
	"context"
)

type BookRepository interface {
	GetBookByID(id string) (Book, error)
	AddBook(ctx context.Context, newBook Book) error
	DeleteBook(ctx context.Context, id string) error
	UpdateBook(ctx context.Context, id string, updatedBook Book) error
}

type BookService interface {
	GetBookByID(id string) (Book, error)
	AddBook(ctx context.Context, newBook Book) error
	DeleteBook(ctx context.Context, id string) error
	UpdateBook(ctx context.Context, id string, updatedBook Book) error
}
