package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/oayomide/xoxo/model"
)

func main() {
	router := mux.NewRouter()
	router.Use(AuthMiddleWare)
	router.HandleFunc("/", handleHomeRoute).Methods("GET")
	router.HandleFunc("/signup", handleSignUp).Methods("POST")
	router.HandleFunc("/api/v1/login", handleLoginRoute).Methods("POST", "OPTIONS", "PUT")
	router.HandleFunc("/me/note/new", HandleCreateNote).Methods("POST")
	router.HandleFunc("/me/note/update", HandleUpdateNote).Methods("POST")
	router.HandleFunc("/me/note/delete", HandleNotesDelete).Methods("DELETE")
	router.HandleFunc("/api/v1/me", HandleMe)

	// update the profile
	router.HandleFunc("/api/v1/user/{id}", UpdateProfile).Methods("POST", "PUT")
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	credentialsOk := handlers.AllowCredentials()

	// get the port
	port := os.Getenv("port")
	if port == "" {
		port = ":13000"
	}

	fmt.Printf("SERVER UP AND RUNNING ON PORT: %s\n", port)
	err := http.ListenAndServe(port, handlers.CORS(headersOk, originsOk, methodsOk, credentialsOk)(router))

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
