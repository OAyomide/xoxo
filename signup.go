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
		res.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	// we're assiging another variable to the model.signup struct cos this will point to
	// another instance of the struct. the one above (var res) is already pointing to the data
	// we got from the endpoint in the frontend
	collection, err := db.GetDBCollection()

	if err != nil {
		fmt.Print("ERROR GETTING THE COLLECTION FROM THE DB.")
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
				res.Error = err.Error()
				json.NewEncoder(w).Encode(res)
				return
			}

			signup.Password = string(hash)
			_, err = collection.InsertOne(context.TODO(), signup)
			if err != nil {
				fmt.Printf("ERROR CREATING USER: %#v\n", err)
				res.Error = err.Error()
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(res)
				return
			}

			res.Data = "NEW USER CREATED!"
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(res)
			return
		}
	}

	res.Data = "USERNAME ALREADY TAKEN"
	w.WriteHeader(http.StatusConflict)
	json.NewEncoder(w).Encode(res)
	return
}
