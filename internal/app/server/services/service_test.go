package services

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/Te8va/APIbook/internal/app/server/domain"
	"github.com/Te8va/APIbook/internal/app/server/domain/mocks"
)

func TestGetBookByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)
	testService := NewBookService(mockRepo)

	book := domain.Book{
		ID:        "123",
		Title:     "Test Book",
		Author:    "Author",
		Year:      2020,
		IsDeleted: false,
	}

	testCases := []struct {
		name     string
		id       string
		wantBook domain.Book
		wantErr  error
	}{
		{
			name:     "valid book ID",
			id:       "123",
			wantBook: book,
			wantErr:  nil,
		},
		{
			name:     "book not found",
			id:       "invalid-id",
			wantBook: book,
			wantErr:  domain.ErrBookNotFound,
		},
		{
			name:     "database error",
			id:       "123",
			wantBook: book,
			wantErr:  domain.ErrReadingDatabase,
		},
		{
			name:     "deleted book",
			id:       "123",
			wantBook: book,
			wantErr:  domain.ErrDeletedBook,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.wantErr != nil {
				mockRepo.EXPECT().GetBookByID(testCase.id).Return(domain.Book{}, testCase.wantErr).Times(1)
			} else {
				mockRepo.EXPECT().GetBookByID(testCase.id).Return(testCase.wantBook, nil).Times(1)
			}

			_, err := testService.GetBookByID(testCase.id)

			if testCase.wantErr != nil {
				require.ErrorIs(t, err, testCase.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAddBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)
	testService := NewBookService(mockRepo)

	book := domain.Book{
		Title:  "Test Book",
		Author: "Author",
		Year:   2020,
	}

	testCases := []struct {
		name    string
		id      string
		newBook domain.Book
		wantErr error
	}{
		{
			name:    "success",
			id:      "123",
			newBook: book,
			wantErr: nil,
		},
		{
			name:    "database error",
			id:      "123",
			newBook: book,
			wantErr: domain.ErrDatabaseOperation,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.wantErr != nil {
				mockRepo.EXPECT().AddBook(gomock.Any(), testCase.newBook).Return("", testCase.wantErr).Times(1)
			} else {
				mockRepo.EXPECT().AddBook(gomock.Any(), testCase.newBook).Return(testCase.id, nil).Times(1)
			}

			curID, err := testService.AddBook(context.Background(), testCase.newBook)

			if testCase.wantErr != nil {
				require.ErrorIs(t, err, testCase.wantErr)
				require.Empty(t, curID)
			} else {
				require.NoError(t, err)
				require.Equal(t, testCase.id, curID)
			}
		})
	}
}

func TestDeleteBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)
	testService := NewBookService(mockRepo)

	testCases := []struct {
		name    string
		id      string
		wantErr error
	}{
		{
			name:    "success",
			id:      "1",
			wantErr: nil,
		},
		{
			name:    "book not found",
			id:      "invalid-id",
			wantErr: domain.ErrBookNotFound,
		},
		{
			name:    "book already deleted",
			id:      "123",
			wantErr: domain.ErrBookNotFound,
		},
		{
			name:    "database error",
			id:      "1",
			wantErr: domain.ErrDatabaseOperation,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.wantErr != nil {
				mockRepo.EXPECT().DeleteBook(gomock.Any(), testCase.id).Return(testCase.wantErr).Times(1)
			} else {
				mockRepo.EXPECT().DeleteBook(gomock.Any(), testCase.id).Return(nil).Times(1)
			}

			err := testService.DeleteBook(context.Background(), testCase.id)

			if testCase.wantErr != nil {
				require.ErrorIs(t, err, testCase.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUpdateBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)
	testService := NewBookService(mockRepo)

	updatedBook := domain.Book{
		Title:  "Updated Book",
		Author: "Updated Author",
		Year:   2000,
	}

	testCases := []struct {
		name        string
		id          string
		updatedBook domain.Book
		wantErr     error
	}{
		{
			name:        "success",
			id:          "1",
			updatedBook: updatedBook,
			wantErr:     nil,
		},
		{
			name:        "book not found",
			id:          "not-found-id",
			updatedBook: updatedBook,
			wantErr:     domain.ErrBookNotFound,
		},
		{
			name:        "book already deleted",
			id:          "deleted-id",
			updatedBook: updatedBook,
			wantErr:     domain.ErrDeletedBook,
		},
		{
			name:        "database error",
			id:          "1",
			updatedBook: updatedBook,
			wantErr:     domain.ErrDatabaseOperation,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.wantErr != nil {
				mockRepo.EXPECT().UpdateBook(gomock.Any(), testCase.id, testCase.updatedBook).Return(testCase.wantErr).Times(1)
			} else {
				mockRepo.EXPECT().UpdateBook(gomock.Any(), testCase.id, testCase.updatedBook).Return(nil).Times(1)
			}

			err := testService.UpdateBook(context.Background(), testCase.id, testCase.updatedBook)

			if testCase.wantErr != nil {
				require.ErrorIs(t, err, testCase.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
