package middleware

type Middlewares struct {
	*Common
}

func New() *Middlewares {
	return &Middlewares{
		Common: NewCommon(),
	}
}
