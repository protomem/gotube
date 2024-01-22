package handler

type User struct {
	*Base
}

func NewUser() *User {
	return &User{NewBase()}
}
