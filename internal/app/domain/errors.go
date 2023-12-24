package domain

import (
	"errors"
)

var (
	ErrDecodingJSON  = errors.New("Error decoding JSON: %v")
	ErrEncodingJSON  = errors.New("Error encoding JSON: %v")
	ErrDeletedBook   = errors.New("Error deleting book")
	ErrReadingFile   = errors.New("Error reading file: %v")
	ErrOpeningFile   = errors.New("Error opening file: %v")
	ErrCreatingFile  = errors.New("Error creating file: %v")
	ErrWritingToFile = errors.New("Error writing to file: %v")
	ErrBookNotFound  = errors.New("Book not found")
)
