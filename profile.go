package main

import (
	"encoding/json"
	"net/http"
)

func HandleMe(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(claims)
}
