package hashing

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var _ Hasher = (*Bcrypt)(nil)

const (
	BcryptMinCost     = bcrypt.MinCost
	BcryptDefaultCost = bcrypt.DefaultCost
	BcryptMaxCost     = bcrypt.MaxCost
)

type Bcrypt struct {
	cost int
}

func NewBcrypt(cost int) *Bcrypt {
	if cost < BcryptMinCost || cost > BcryptMaxCost {
		cost = BcryptDefaultCost
	}

	return &Bcrypt{cost: cost}
}

func (b *Bcrypt) Generate(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", fmt.Errorf("hashing:Bcrypt.Generate: %w", err)
	}

	return string(hash), nil
}

func (b *Bcrypt) Compare(password string, hashed string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return fmt.Errorf("hashing:Bcrypt.Compare: %w", ErrWrongPassword)
		}

		return fmt.Errorf("hashing:Bcrypt.Compare: %w", err)
	}

	return nil
}
