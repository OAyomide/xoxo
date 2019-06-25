package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/mongodb/mongo-go-driver/bson"
	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"
	"github.com/oayomide/xoxo/db"
	"github.com/oayomide/xoxo/model"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", handleHomeRoute).Methods("GET")
	router.HandleFunc("/signup", handleSignUp).Methods("POST")
	router.HandleFunc("/login", handleLoginRoute).Methods("POST")
	router.HandleFunc("/me", handleProfileRoute).Methods("GET")
	// http.HandleFunc("/", handleHomeRoute)
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

func handleSignUp(w http.ResponseWriter, r *http.Request) {
	// set the header to allow get and post requests
	w.Header().Set("Content-Type", "application/json")
	var signup model.Signup
	body, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(body, &signup)
	var res model.Response

	if err != nil {
		fmt.Printf("ERROR PARSING JSON AND SAVING TO POINTER TO SIGNUP STRUCT %#v", err)
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	// we're assiging another variable to the model.signup struct cos this will point to
	// another instance of the struct. the one above (var res) is already pointing to the data
	// we got from the endpoint in the frontend
	collection, err := db.GetDBCollection()

	if err != nil {
		fmt.Printf("ERROR GETTING THE COLLECTION FROM THE DB. %#v", err)
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	var result model.Signup
	err = collection.FindOne(context.TODO(), bson.D{{"username", signup.Username}}).Decode(&result)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			hash, err := bcrypt.GenerateFromPassword([]byte(signup.Password), 6)

			if err != nil {
				fmt.Printf("ERROR HASHING USER PASSWORD USING BCRYPT: %#v", err)
				res.Error = err.Error()
				json.NewEncoder(w).Encode(res)
				return
			}

			signup.Password = string(hash)
			newlyRegisterdUser, err := collection.InsertOne(context.TODO(), signup)

			if err != nil {
				fmt.Printf("ERROR CREATING USER: %#v", err)
				res.Error = err.Error()
				json.NewEncoder(w).Encode(res)
				return
			}

			fmt.Printf("NEWLY REGISTERED USER IS: %#v", newlyRegisterdUser)
			res.Data = "NEW USER CREATED!"
			json.NewEncoder(w).Encode(res)
			return
		}
	}

	res.Data = "USERNAME ALREADY TAKEN"
	json.NewEncoder(w).Encode(res)
	return
}

func handleLoginRoute(w http.ResponseWriter, r *http.Request) {
	// first set the header to application/json
	w.Header().Set("Content-Type", "application/json")
	var user model.User

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)

	fmt.Printf("MARSHALLED:: %#v", user)
	if err != nil {
		fmt.Printf("ERROR DECODING USER STRUCT TO JSON: %#v", err)
		panic(err)
	}

	collection, err := db.GetDBCollection()
	if err != nil {
		fmt.Printf("ERROR GETTING DB COLLECTION: %#v", err)
		panic(err)
	}
	var result model.User
	var res model.Response

	err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result) // saving the result into a pointer to another instance of model.User
	if err != nil {
		fmt.Print("ERROR GETTING USERNAME FROM THE DB: %#v", err)
		res.Error = "USERNAME OR PASSWORD NOT VALID"
		json.NewEncoder(w).Encode(res)
		return
	}
	fmt.Printf("USER PASSWORD FROM THE API IS: %s\n", result.Password)
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
	if err != nil {
		fmt.Printf("USER PASSWORD INVALID: %#v", err)
		res.Error = "USERNAME OR PASSWORD NOT VALID"
		json.NewEncoder(w).Encode(res)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": result.Username,
		"phone":    result.Phone,
		"email":    result.Email,
	})

	tokenString, err := token.SignedString([]byte("xoxo")) //TODO: Move this to a config file/struct
	if err != nil {
		fmt.Printf("ERROR SIGNING STRING:: %#v", err)
		json.NewEncoder(w).Encode(res)
		return
	}

	result.Token = tokenString
	result.Password = ""
	json.NewEncoder(w).Encode(result) // encoding result because thats what we want to return. it contains the jwt key
	return
}

func handleProfileRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, fmt.Errorf("UNEXPECTED TOKEN SIGNING METHOD")
		}

		return []byte("xoxo"), nil
	})

	var result model.User
	var res model.Response

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		result.Username = claims["username"].(string)
		result.Email = claims["email"].(string)
		result.Phone = claims["phone"].(string)

		json.NewEncoder(w).Encode(result)
		return
	} else {
		res.Error = err.Error()
		fmt.Printf("ERROR GETTING USER PROFILE HERE: %#v", err)
		json.NewEncoder(w).Encode(res)
		return
	}
	// body, _ := ioutil.ReadAll(r.Body)
}
