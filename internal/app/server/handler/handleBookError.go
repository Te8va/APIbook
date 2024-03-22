package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Te8va/APIbook/internal/app/server/domain"
	logging "github.com/Te8va/APIbook/internal/pkg/logger"
)

func handleBookError(w http.ResponseWriter, err error) {
	var statusCode int

	if errors.Is(err, domain.ErrReadingFile) || errors.Is(err, domain.ErrDecodingJSON) || errors.Is(err, domain.ErrEncodingJSON) || errors.Is(err, domain.ErrWritingToFile) {
		statusCode = http.StatusInternalServerError
	} else if errors.Is(err, domain.ErrDeletedBook) || errors.Is(err, domain.ErrBookNotFound) {
		statusCode = http.StatusNotFound
	} else {
		statusCode = http.StatusInternalServerError
		err = fmt.Errorf("internal server error: %w", err)
	}

	logging.Logger().Error(err.Error(), err)
	http.Error(w, err.Error(), statusCode)
}
