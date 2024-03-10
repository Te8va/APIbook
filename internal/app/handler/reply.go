package handler

import (
	"encoding/json"
	"net/http"

	logging "github.com/Te8va/APIbook/internal/pkg/logger"
)

func reply(w http.ResponseWriter, message interface{}, statusCode int) {

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
