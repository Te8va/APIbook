package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func reply(w http.ResponseWriter, message interface{}, statusCode int) {

	response, err := json.Marshal(message)
	if err != nil {
		handleBookError(w, fmt.Errorf("Error encoding JSON: %v", err))
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	// TODO: do it externally, write status code after write, check write error
	w.WriteHeader(statusCode)
	w.Write(response)

}
