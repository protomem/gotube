package entity

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
	Entity string
	Err    error
}

func NewError(entity string, err error) Error {
	return Error{
		Entity: entity,
		Err:    err,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Entity, e.Err.Error())
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

func IsError(err error, target Error) bool {
	var e Error
	if !errors.As(err, &e) {
		return false
	}
	return strings.EqualFold(e.Entity, target.Entity) && errors.Is(e.Err, target.Err)
}
