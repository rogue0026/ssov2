package models

type User struct {
	ID           int64
	Login        string
	PasswordHash []byte
	Email        string
}
