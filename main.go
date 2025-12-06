package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	// Just one very simple handler for test
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		response := struct {
			Status string `json:"status"`
		}{Status: "OK"}

		json.NewEncoder(w).Encode(response)
	})
	_ = http.ListenAndServe(":8000", nil)
}
