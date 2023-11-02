package bcrypt

import (
	"errors"
	"fmt"

	"github.com/protomem/gotube/internal/hashing"
	"golang.org/x/crypto/bcrypt"
)

const (
	MinCost     = bcrypt.MinCost
	DefaultCost = bcrypt.DefaultCost
	MaxCost     = bcrypt.MaxCost
)

var _ hashing.Hasher = (*Hasher)(nil)

type Hasher struct {
	cost int
}

func New(cost int) *Hasher {
	return &Hasher{
		cost: cost,
	}
}

func (h *Hasher) Generate(password string) (string, error) {
	const op = "bcrypt.Generate"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return string(hash), nil
}

func (h *Hasher) Compare(password string, hash string) error {
	const op = "bcrypt.Compare"

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return fmt.Errorf("%s: %w", op, hashing.ErrWrongPassword)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
