package model

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
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone,omitempty"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type Response struct {
	Error string `json:"error"`
	Data  string `json:"data,omitempty"`
}
