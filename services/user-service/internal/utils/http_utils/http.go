package http_utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func JsonResponse(w http.ResponseWriter, statusCode int, output interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(output)
	if err != nil {
		fmt.Printf("error while encoding response: %v\n", err)
	}
}

func ErrorInvalidRequestBody(w http.ResponseWriter, err error) {
	JsonError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %s", err.Error()))
}

func SuccessJsonResponse(w http.ResponseWriter, output interface{}) {
	JsonResponse(w, http.StatusOK, output)
}

func JsonError(w http.ResponseWriter, statusCode int, err error) {
	JsonResponse(w, statusCode, &ErrorResponse{Error: err.Error()})
}

func JsonErrorNotFound(w http.ResponseWriter) {
	JsonError(w, http.StatusNotFound, fmt.Errorf("not found"))
}

func JsonInternalServerError(w http.ResponseWriter) {
	JsonError(w, http.StatusInternalServerError, fmt.Errorf("internal error"))
}
