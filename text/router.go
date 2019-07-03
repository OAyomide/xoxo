package text

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mongodb/mongo-go-driver/mongo/options"

	"github.com/mongodb/mongo-go-driver/bson"

	"github.com/dgrijalva/jwt-go"

	"github.com/oayomide/xoxo/db"
	"github.com/oayomide/xoxo/model"
)

type HandleCopyResponse struct {
	Username string `json:"username"`
	Note     string `json:"note"`
	Time     string `json:"timestamp"`
}

func HandleTextCopy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	var text Text
	var res model.Response
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &text)

	if err != nil {
		fmt.Printf("ERROR PARSING JSON INTO TEXT RECEIVER")
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	collection, _ := db.GetCollection("user")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("UNEXPECTED JWT TOKEN")
		}

		return []byte("xoxo"), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)

	var result Text
	if ok && token.Valid {
		// uid, _ := primitive.ObjectIDFromHex(claims["_id"].(string))
		uid := claims["_id"].(string)

		if uid == "" {
			fmt.Println("USER UID NOT FOUND.. USER IS NOT VALID")
			res.Error = "user not found"
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(res)
			return
		}
		err := collection.FindOne(context.TODO(), bson.D{{"id", uid}}).Decode(&result)

		if err != nil {
			textCollection, _ := db.GetCollection("note")
			err = textCollection.FindOne(context.TODO(), bson.D{{"user", uid}}).Decode(&result)

			if err != nil {
				fmt.Println("USER DOESNT HAVE ANY NOTE CREATED. . .GOING TO CREATE FOR USER")
				_, err := textCollection.InsertOne(context.TODO(), bson.D{{"user", uid}, {"note", text.Note}, {"timestamp", text.Timestamp}})
				if err != nil {
					res.Error = err.Error()
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(res)
					return
				}

				res.Data = "NOTE CREATED FOR NEW USER"
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(res)
				return
			}

			// _, err := textCollection.InsertOne(context.TODO(), bson.D{{"user", uid}, {"note", text.Note}, {"timestamp", text.Timestamp}})
			_, err := textCollection.UpdateOne(context.TODO(), bson.D{{"user", uid}}, bson.D{{"$set", bson.D{{"note", text.Note}, {"timestamp", text.Timestamp}}}}, options.Update().SetUpsert(true))

			if err != nil {
				fmt.Println("ERROR CREATING AND UPDATING NEW USER DOCUMENT INTO THE DB")
				res.Error = err.Error()
				json.NewEncoder(w).Encode(res)
				return
			}

			var hnresponse HandleCopyResponse
			hnresponse.Note = text.Note
			hnresponse.Time = string(text.Timestamp)
			res.Data = "Created note for user!"
			json.NewEncoder(w).Encode(hnresponse)
			return
		}
	} else {
		res.Error = err.Error()
		fmt.Println("TOKEN IS NOT VALID")
		json.NewEncoder(w).Encode(res)
	}
}

// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 	_, ok := token.Method.(*jwt.SigningMethodHMAC)

// 	if !ok {
// 		return nil, fmt.Errorf("UNEXPECTED TOKEN")
// 	}

// 	return []byte("xoxo"), nil
// })

// var text Text
// var res model.Response

// collection, err := db.GetCollection("users")
// body, _ := ioutil.ReadAll(r.Body)

// err = json.Unmarshal(body, &text)
// if err != nil {
// 	fmt.Printf("ERROR PARSING THE JSON INTO THE text POINTER")
// 	res.Error = err.Error()
// 	json.NewEncoder(w).Encode(res)
// 	return
// }
// _, ok := token.Claims.(jwt.MapClaims)

// var result Text
// if ok && token.Valid {
// 	// first look for the user in the DB
// 	err = collection.FindOne(context.TODO(), bson.D{{"userid": uid}}).Decode(&result)
// 	// then enter into the DB
// 	if err != nil {

// 	}
// }
