package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/oayomide/xoxo/db"
	"github.com/oayomide/xoxo/model"
	"golang.org/x/crypto/bcrypt"
)

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
	var lresponse model.LoginResponse
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
		"_id":      result.ID,
	})

	tokenString, err := token.SignedString([]byte("xoxo")) //TODO: Move this to a config file/struct
	if err != nil {
		fmt.Printf("ERROR SIGNING STRING:: %#v", err)
		json.NewEncoder(w).Encode(res)
		return
	}

	lresponse.Email = result.Email
	lresponse.ID = result.ID.Hex()
	lresponse.Phone = result.Phone
	lresponse.Username = result.Username
	lresponse.Token = tokenString

	json.NewEncoder(w).Encode(lresponse) // encoding result because thats what we want to return. it contains the jwt key
	return
}
