package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/Te8va/APIbook/internal/app/logging"
)

func reply(w http.ResponseWriter, message interface{}, statusCode int) {

	response, err := json.Marshal(message)
	if err != nil {
		logging.Logger.Error("Error encoding JSON", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_, err = w.Write(response)
	if err != nil {
		logging.Logger.Error("Error writing response", zap.Error(err))
		return
	}

}
