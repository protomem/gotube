package model

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type Error struct {
	Model string
	Err   error
}

func NewError(model string, err error) Error {
	return Error{
		Model: model,
		Err:   err,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Model, e.Err.Error())
}

func (e Error) Is(target error) bool {
	return errors.Is(e.Err, target)
}

func (e Error) As(target any) bool {
	_, ok := target.(Error)
	if !ok {
		return errors.As(e.Err, target)
	}
	return true
}

func IsModelError(err error, target Error) bool {
	var e Error
	if !errors.As(err, &e) {
		return false
	}
	return strings.EqualFold(e.Model, target.Model) && errors.Is(e.Err, target.Err)
}