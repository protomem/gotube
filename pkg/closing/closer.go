package closing

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type Func func(context.Context) error

type Closer struct {
	mux sync.Mutex
	fns []Func
}

func New() *Closer {
	return &Closer{
		fns: make([]Func, 0),
	}
}

func (c *Closer) Add(fn Func) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.fns = append(c.fns, fn)
}

func (c *Closer) Close(ctx context.Context) error {
	const op = "closer.Close"

	c.mux.Lock()
	defer c.mux.Unlock()

	var (
		msgs     = make([]string, 0, len(c.fns))
		complete = make(chan struct{}, 1)
	)

	go func() {
		for _, fn := range c.fns {
			if err := fn(ctx); err != nil {
				msgs = append(msgs, fmt.Sprintf("%v", err))
			}
		}

		complete <- struct{}{}
	}()

	select {
	case <-complete:
		break
	case <-ctx.Done():
		return fmt.Errorf("%s: shutdown canceled: %w", op, ctx.Err())
	}

	if len(msgs) > 0 {
		return fmt.Errorf(
			"%s: shutdown finished with error(s): %v",
			op, strings.Join(msgs, "; "),
		)
	}

	return nil
}
