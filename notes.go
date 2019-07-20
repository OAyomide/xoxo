package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"github.com/mongodb/mongo-go-driver/mongo/options"

	"github.com/mongodb/mongo-go-driver/bson"

	"github.com/oayomide/xoxo/db"
	"github.com/oayomide/xoxo/model"
)

func HandleCreateNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	uid := r.Context().Value("claims")
	var text model.Text
	var res model.Response
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &text)

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	ntTSave := []model.Text{model.Text{Title: text.Title, Note: text.Note, Timestamp: time.Now().String(), ID: primitive.NewObjectID()}}
	collection, _ := db.GetCollection("users")
	// we want to find user selection in notes collection
	notesCollection, _ := db.GetCollection("notes")
	var result model.UserNotes
	var user model.User
	id, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", uid))
	err = collection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&user)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			fmt.Println("USER NOT FOUND FROM THE DB")
			res.Error = "user doesnt exist. cannot create note"
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(res)
			return
		}
	}

	err = notesCollection.FindOne(context.TODO(), bson.D{{"user", id.Hex()}}).Decode(&result)

	// user doesnt exist. has no note created at all
	if result.User == "" {
		_, err = notesCollection.InsertOne(context.TODO(), bson.D{{"user", id.Hex()}, {"notes", ntTSave}})

		if err != nil {
			res.Data = "not created"
			res.Error = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(res)
			return
		}

		fmt.Println("SUCCESSFULLY CREATED NEW NOTE FOR THE USER")
		res.Data = "created"
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res)
		return
	} else {
		var nt model.UserNotes
		// user has a note created.. then look if the user has a note with the name already
		_ = notesCollection.FindOne(context.TODO(), bson.D{{"user", id.Hex()}, {"notes.title", text.Title}}).Decode(&nt)
		// the notes that has the same title as the note we want to create exists. i.e it doesnt return an empty array
		if len(nt.Notes) > 0 {
			res.Error = "note already exists"
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(res)
			return
		}

		_, inserError := notesCollection.InsertOne(context.TODO(), bson.D{{"user", id.Hex()}, {"notes", ntTSave}})

		if inserError != nil {
			res.Data = "error creating"
			res.Error = inserError.Error()
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(res)
			return
		}
		res.Data = "created"
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res)
		return
	}
}

func HandleUpdateNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	uid := r.Context().Value("claims")
	id, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", uid))
	var text model.Text
	var res model.Response
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &text)

	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}
	var result model.Text
	notesCollection, _ := db.GetCollection("notes")
	err = notesCollection.FindOne(context.TODO(), bson.D{{"user", id.Hex()}, {"notes.title", text.Title}}).Decode(&result)

	if err != nil {
		fmt.Println("COULNDT FIND THE NOTE WITH THAT TITLE")
		res.Error = err.Error()
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(res)
		return
	}

	_, err = notesCollection.UpdateOne(context.TODO(), bson.D{{"user", id.Hex()}, {"notes.title", text.Title}}, bson.D{{"$set", bson.D{{"notes.$.note", text.Note}, {"notes.$.timestamp", time.Now().String()}}}}, options.Update().SetUpsert(true))

	if err != nil {
		fmt.Println("ERROR UPDATING DOCUMENT IN THE DB")
		res.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	var hnresponse model.HandleCopyResponse
	hnresponse.Note = text.Note
	hnresponse.Time = string(time.Now().String())
	res.Data = "Created note for user!"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(hnresponse)
	return
}

func HandleNotesDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	uid := r.Context().Value("claims")
	id, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", uid))

	// then get the id from the db
	var note model.Text
	var res model.Response

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &note)
	notesCollection, _ := db.GetCollection("notes")
	if err != nil {
		fmt.Println("COULNDT FIND THE NOTE WITH THAT TITLE")
		res.Error = err.Error()
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(res)
		return
	}

	_, err = notesCollection.DeleteOne(context.TODO(), bson.D{{"user", id.Hex()}, {"notes.title", note.Title}})
	if err != nil {
		fmt.Println("ERROR DELETING NOTE")
		res.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	fmt.Println("DELETED NOTE")
	res.Data = "note deleted"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
	return
}
