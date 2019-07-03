package model

import (
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type HelloResponse struct {
	Message string `json:"message"`
}

type Signup struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone,omitempty"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type User struct {
	Username string             `json:"username"`
	Password string             `json:"password"`
	Phone    string             `json:"phone,omitempty"`
	Email    string             `json:"email"`
	Token    string             `json:"token"`
	ID       primitive.ObjectID `bson:"_id"`
}

type Notes struct {
	User      string `json:"id"`
	Text      string `json:"note"`
	Timestamp string `json:"timestamp"`
}

type Response struct {
	Error string `json:"error"`
	Data  string `json:"data,omitempty"`
}
