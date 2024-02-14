package handler

type Handlers struct {
	*Common
}

func New() *Handlers {
	return &Handlers{
		Common: NewCommon(),
	}
}
