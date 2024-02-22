package rest

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	response := ErrorResponse{Message: "Internal Server Error"}
	json.NewEncoder(w).Encode(response)
}

func NotFoundErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	response := ErrorResponse{Message: "Not Found"}
	json.NewEncoder(w).Encode(response)
}

func BadRequestErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	response := ErrorResponse{Message: "Bad Request"}
	json.NewEncoder(w).Encode(response)
}

func UnauthorizedErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	response := ErrorResponse{Message: "Unauthorized"}
	json.NewEncoder(w).Encode(response)
}
