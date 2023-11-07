package requestid

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const (
	Header = "X-Request-Id"
	LogKey = "requestId"
)

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rid := r.Header.Get(Header)
			if rid == "" {
				rid = Generate()
				r.Header.Set(Header, rid)
			}

			ctx := Inject(r.Context(), rid)
			r = r.WithContext(ctx)

			w.Header().Set(Header, rid)

			next.ServeHTTP(w, r)
		})
	}
}

func Empty() string {
	return uuid.Nil.String()
}

func Generate() string {
	rid, err := uuid.NewRandom()
	if err != nil {
		return Empty()
	}

	return rid.String()
}

type ctxKey struct{}

func Inject(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, ctxKey{}, rid)
}

func Extract(ctx context.Context) string {
	rid, ok := ctx.Value(ctxKey{}).(string)
	if !ok {
		return Empty()
	}

	return rid
}
