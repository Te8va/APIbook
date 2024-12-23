package domain

import (
	"errors"
)

var (
	ErrDecodingJSON      = errors.New("error decoding JSON")
	ErrEncodingJSON      = errors.New("error encoding JSON")
	ErrDeletedBook       = errors.New("error deleting book")
	ErrReadingFile       = errors.New("error reading file")
	ErrOpeningFile       = errors.New("error opening file")
	ErrCreatingFile      = errors.New("error creating file")
	ErrWritingToFile     = errors.New("error writing to file")
	ErrBookNotFound      = errors.New("book not found")
	ErrReadingDatabase   = errors.New("error reading database")
	ErrDatabaseOperation = errors.New("error database operation")
	ErrCreatingDirectory = errors.New("error creating directory")
)
