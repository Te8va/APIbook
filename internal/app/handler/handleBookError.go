package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Te8va/APIbook/internal/app/domain"
)

func handleBookError(w http.ResponseWriter, err error) {
	var statusCode int
	var errorMessage string

	switch {
	case errors.Is(err, domain.ErrReadingFile):
		statusCode = http.StatusInternalServerError
		errorMessage = fmt.Sprintf("Error reading file: %v", err)
	case errors.Is(err, domain.ErrDecodingJSON):
		statusCode = http.StatusInternalServerError
		errorMessage = fmt.Sprintf("Error decoding JSON: %v", err)
	case errors.Is(err, domain.ErrDeletedBook):
		statusCode = http.StatusNotFound
		errorMessage = fmt.Sprintf("Error deleting book: %v", err)
	case errors.Is(err, domain.ErrEncodingJSON):
		statusCode = http.StatusInternalServerError
		errorMessage = fmt.Sprintf("Error encoding JSON: %v", err)
	case errors.Is(err, domain.ErrWritingToFile):
		statusCode = http.StatusInternalServerError
		errorMessage = fmt.Sprintf("Error writing to file: %v", err)
	case errors.Is(err, domain.ErrBookNotFound):
		statusCode = http.StatusNotFound
		errorMessage = fmt.Sprintf("Book not found: %v", err)
	default:
		statusCode = http.StatusInternalServerError
		errorMessage = fmt.Sprintf("Internal server error: %v", err)
	}

	http.Error(w, errorMessage, statusCode)
}
