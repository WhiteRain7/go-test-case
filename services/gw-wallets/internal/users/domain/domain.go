package domain

type Hashed = string

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Credentials struct {
	Username string
	Password Hashed
}
