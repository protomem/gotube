package bcrypt

import (
	"errors"
	"fmt"

	"github.com/protomem/gotube/pkg/hashing"
	"golang.org/x/crypto/bcrypt"
)

var _ hashing.Hasher = (*Hasher)(nil)

const (
	MinCost     = bcrypt.MinCost
	MaxCost     = bcrypt.MaxCost
	DefaultCost = bcrypt.DefaultCost
)

type Hasher struct {
	cost int
}

func New(cost int) *Hasher {
	return &Hasher{
		cost: cost,
	}
}

func (h *Hasher) Generate(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("hasher.Generate: %w", err)
	}

	return string(hash), nil
}

func (h *Hasher) Verify(password, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return fmt.Errorf("hasher.Verify: %w", hashing.ErrWrongPassword)
		}
		return fmt.Errorf("hasher.Verify: %w", err)
	}
	return nil
}
