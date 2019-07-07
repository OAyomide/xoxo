package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/oayomide/xoxo/db"
	"github.com/oayomide/xoxo/model"
	"golang.org/x/crypto/bcrypt"
)

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
