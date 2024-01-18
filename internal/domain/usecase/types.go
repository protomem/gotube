package usecase

import "context"

type Usecase[I any, O any] interface {
	Invoke(ctx context.Context, input I) (output O, err error)
}

type UsecaseFunc[I any, O any] func(ctx context.Context, input I) (output O, err error)

func (fn UsecaseFunc[I, O]) Invoke(ctx context.Context, input I) (O, error) {
	return fn(ctx, input)
}

type FindOptions struct {
	Limit  uint64
	Offset uint64
}

type void struct{}
