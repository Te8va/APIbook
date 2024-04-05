package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Te8va/APIbook/internal/app/server/domain"
)

type postgres struct {
	*pgxpool.Pool
}

func NewPostgresBookRepository(pg *pgxpool.Pool) *postgres {
	return &postgres{Pool: pg}
}

func (p *postgres) GetBookByID(id string) (domain.Book, error) {
	var book domain.Book

	err := p.QueryRow(context.Background(), "SELECT id, title, author, year, status FROM books WHERE id = $1", id).Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.Status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Book{}, domain.ErrBookNotFound
		}
		return domain.Book{}, domain.ErrReadingDatabase
	}

	if book.Status == "deleted" {
		return domain.Book{}, domain.ErrDeletedBook
	}

	return book, nil
}

func (p *postgres) AddBook(ctx context.Context, newBook domain.Book) (string, error) {
	err := p.UseTransaction(ctx, func(tx pgx.Tx) error {

		newUUID := uuid.New()
		newBook.ID = newUUID.String()

		_, err := tx.Exec(ctx, "INSERT INTO books (id, title, author, year, status) VALUES ($1, $2, $3, $4, $5)", newBook.ID, newBook.Title, newBook.Author, newBook.Year, "")
		if err != nil {
			return domain.ErrDatabaseOperation
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return newBook.ID, nil
}

func (p *postgres) UpdateBook(ctx context.Context, id string, updatedBook domain.Book) error {
	var book domain.Book

	err := p.QueryRow(ctx, "SELECT status FROM books WHERE id = $1", id).Scan(&book.Status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ErrBookNotFound
		}
		return domain.ErrReadingDatabase
	}

	if book.Status == "deleted" {
		return domain.ErrDeletedBook
	}

	return p.UseTransaction(ctx, func(tx pgx.Tx) error {

		tag, err := tx.Exec(ctx, "UPDATE books SET title = $1, author = $2, year= $3 WHERE id = $4", updatedBook.Title, updatedBook.Author, updatedBook.Year, id)
		if err != nil {
			return domain.ErrDatabaseOperation
		}
		if tag.RowsAffected() == 0 {
			return domain.ErrBookNotFound
		}

		return nil
	})
}

func (p *postgres) DeleteBook(ctx context.Context, id string) error {
	var book domain.Book

	err := p.QueryRow(ctx, "SELECT status FROM books WHERE id = $1", id).Scan(&book.Status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ErrBookNotFound
		}
		return domain.ErrReadingDatabase
	}

	if book.Status == "deleted" {
		return domain.ErrDeletedBook
	}

	return p.UseTransaction(ctx, func(tx pgx.Tx) error {

		_, err := tx.Exec(ctx, "UPDATE books SET status = 'deleted' WHERE id = $1", id)
		if err != nil {
			return domain.ErrDatabaseOperation
		}

		return nil
	})
}

func NewPgxpool(dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func (p *postgres) UseTransaction(ctx context.Context, txFunc func(pgx.Tx) error) error {
	tx, err := p.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			return
		}
	}()

	err = txFunc(tx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
