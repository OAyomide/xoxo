package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"github.com/mongodb/mongo-go-driver/bson"

	"github.com/gorilla/mux"
	"github.com/oayomide/xoxo/db"
	"github.com/oayomide/xoxo/model"
)

// HandleMe returns the current user signed in
func HandleMe(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims")
	fmt.Printf("%v", claims)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(claims)
}

// UpdateProfile updates the profile for the user
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	userID := params["id"]

	fmtUserID, _ := primitive.ObjectIDFromHex(userID)
	var userProfileM model.ProfileResponse
	var errorResponse model.Response
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &userProfileM)

	if err != nil {
		errorResponse.Error = err.Error()
		json.NewEncoder(w).Encode(errorResponse)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// look for the user in the db
	userCollection, _ := db.GetCollection("users")

	var userM model.User
	err = userCollection.FindOne(context.TODO(), bson.D{{"_id", fmtUserID}}).Decode(&userM)

	if err != nil {
		if err.Error() == "mongo: no document in result" {
			fmt.Println("User with the id doesnt exist")
			errorResponse.Error = "User not found"
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
	}

	// then we want to replace the name of the user. but first, we want to see if the
	// username already exists
	var userD model.User
	err = userCollection.FindOne(context.TODO(), bson.D{{"username", userProfileM.Username}}).Decode(&userD)

	fmt.Printf("USER DATA FROM ENDPOING IS:: %v\n\n", userD)
	if userD.Username == "" {
		// username not taken, we want to update the user's name here
		_, err = userCollection.UpdateOne(context.TODO(), bson.D{{"username", userM.Username}}, bson.D{{"$set", bson.D{{"username", userProfileM.Username}, {"phone", userProfileM.Phone}, {"email", userProfileM.Email}}}})

		if err != nil {
			errorResponse.Error = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		errorResponse.Data = "user profile updated"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(errorResponse)
		return
		// if err.Error() == "mongo: no document in result" {

		// }
		// errorResponse.Error = err.Error()
		// w.WriteHeader(http.StatusInternalServerError)
		// json.NewEncoder(w).Encode(errorResponse)
		// return
	}

	// if userD.Username != "" {
	// 	errorResponse.Error = "Username already taken"
	// 	w.WriteHeader(http.StatusConflict)
	// 	json.NewEncoder(w).Encode(errorResponse)
	// 	return
	// }
}
