package port

import "context"

type Void struct{}

type FindOptions struct {
	Limit  uint64
	Offset uint64
}

type Tuple[A, B any] struct {
	A A
	B B
}

func NewTuple[A, B any](a A, b B) Tuple[A, B] {
	return Tuple[A, B]{A: a, B: b}
}

type Usecase[I, O any] interface {
	Invoke(ctx context.Context, input I) (output O, err error)
}

type UsecaseFunc[I, O any] func(ctx context.Context, input I) (output O, err error)

func (fn UsecaseFunc[I, O]) Invoke(ctx context.Context, input I) (O, error) {
	return fn(ctx, input)
}
