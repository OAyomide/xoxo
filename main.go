package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/oayomide/xoxo/model"
)

func main() {
	router := mux.NewRouter()
	router.Use(AuthMiddleWare)
	router.HandleFunc("/", handleHomeRoute).Methods("GET")
	router.HandleFunc("/signup", handleSignUp).Methods("POST")
	router.HandleFunc("/login", handleLoginRoute).Methods("POST")
	router.HandleFunc("/me/note", HandleTextCopy).Methods("POST")
	router.HandleFunc("/me", HandleMe)

	// get the port
	port := os.Getenv("port")
	if port == "" {
		port = ":13000"
	}

	fmt.Printf("SERVER UP AND RUNNING ON PORT: %s\n", port)
	err := http.ListenAndServe(port, router)

	if err != nil {
		panic(err)
	}
}

func handleHomeRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	m := model.HelloResponse{Message: "Hello there!"}
	fmt.Print("REQUEST MADE TO THE HOME ROUTE")
	json.NewEncoder(w).Encode(m)
	return
}
