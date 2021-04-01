package helpers

import (
	"encoding/json"
	"movies-app/structs"
	"net/http"
)

//OK - sends response in json format
func OK(response interface{}, w http.ResponseWriter) {
	data, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// GetError : This is helper function to prepare error model.
func GetError(err error, w http.ResponseWriter) {
	var response = structs.ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}
	message, _ := json.Marshal(response)
	w.WriteHeader(response.StatusCode)
	w.Write(message)
}
