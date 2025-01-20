package model

type User struct {
	Id           string `json:"id" bun:"id,pk"`
	Firstname    string `json:"first_name" bun:"first_name"`
	Lastname     string `json:"last_name" bun:"last_name"`
	Email        string `json:"email" bun:"email"`
	PasswordHash string `json:"-" bun:"password_hash"`
}
