package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func HandleMe(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims")
	fmt.Printf("%v", claims)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(claims)
}
