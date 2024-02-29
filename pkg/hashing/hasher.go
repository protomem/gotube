package hashing

import "fmt"

var ErrWrongPassword = fmt.Errorf("wrong password")

type Hasher interface {
	Generate(password string) (string, error)
	Verify(password, hash string) error
}
