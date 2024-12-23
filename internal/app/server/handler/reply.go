package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Te8va/APIbook/internal/app/server/domain"
	logging "github.com/Te8va/APIbook/internal/pkg/logger"
)

func reply(w http.ResponseWriter, message interface{}, statusCode int, err error) {
	if err != nil {
		if !errors.Is(err, domain.ErrDeletedBook) && !errors.Is(err, domain.ErrBookNotFound) {
			statusCode = http.StatusInternalServerError
		} else {
			statusCode = http.StatusNotFound
		}

		logging.Logger().Error(err.Error())
		http.Error(w, err.Error(), statusCode)
		return
	}

	response, err := json.Marshal(message)
	if err != nil {
		logging.Logger().Error("Error encoding JSON", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_, err = w.Write(response)
	if err != nil {
		logging.Logger().Error("Error writing response", err)
		return
	}

}
