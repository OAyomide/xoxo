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
	Note      string `json:"note"`
	Title     string `json:"title"`
	Timestamp string `json:"timestamp"`
}

type UserNotes struct {
	ID    primitive.ObjectID `bson:"_id"`
	User  string             `json:"user"`
	Notes []Notes            `json:"notes"`
}
type Response struct {
	Error string `json:"error"`
	Data  string `json:"data,omitempty"`
}

type LoginResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"jwt_token,omitempty"`
	Phone    string `json:"phone"`
	ID       string `json:"id"`
}

type ProfileResponse struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	ID       string `json:"id"`
}

type Text struct {
	Note      string             `json:"note"`
	Timestamp string             `json:"timestamp"`
	Title     string             `json:"title"`
	ID        primitive.ObjectID `bson:"_id"`
}

type HandleCopyResponse struct {
	Username string `json:"username"`
	Note     string `json:"note"`
	Time     string `json:"timestamp"`
}
