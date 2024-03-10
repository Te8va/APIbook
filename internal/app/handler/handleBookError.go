package handler

import (
	"fmt"
	"net/http"

	"github.com/Te8va/APIbook/internal/app/domain"
	logging "github.com/Te8va/APIbook/internal/pkg/logger"
)

func handleBookError(w http.ResponseWriter, err error) {
	var statusCode int
	var errorMessage string

	switch err {
	case domain.ErrReadingFile:
		statusCode = http.StatusInternalServerError
		errorMessage = fmt.Sprintf("Error reading file: %v", err)
	case domain.ErrDecodingJSON:
		statusCode = http.StatusInternalServerError
		errorMessage = fmt.Sprintf("Error decoding JSON: %v", err)
	case domain.ErrDeletedBook:
		statusCode = http.StatusNotFound
		errorMessage = fmt.Sprintf("Error deleting book: %v", err)
	case domain.ErrEncodingJSON:
		statusCode = http.StatusInternalServerError
		errorMessage = fmt.Sprintf("Error encoding JSON: %v", err)
	case domain.ErrWritingToFile:
		statusCode = http.StatusInternalServerError
		errorMessage = fmt.Sprintf("Error writing to file: %v", err)
	case domain.ErrBookNotFound:
		statusCode = http.StatusNotFound
		errorMessage = fmt.Sprintf("Book not found: %v", err)
	default:
		statusCode = http.StatusInternalServerError
		errorMessage = fmt.Sprintf("Internal server error: %v", err)
	}

	logging.Logger().Error(errorMessage, err)
	http.Error(w, errorMessage, statusCode)
}
