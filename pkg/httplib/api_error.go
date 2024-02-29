package httplib

type APIError struct {
	Code     int    `json:"-"`
	Message  string `json:"message"`
	Internal error  `json:"-"`
}

func NewAPIError(code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

func (e *APIError) WithInternal(err error) *APIError {
	e.Internal = err
	return &APIError{
		Code:     e.Code,
		Message:  e.Message,
		Internal: err,
	}
}

func (e *APIError) Error() string {
	if e.Internal == nil {
		return e.Message
	}
	return e.Message + ": " + e.Internal.Error()
}

func (e *APIError) As(target any) bool {
	if err, ok := target.(*APIError); ok {
		*err = *e
		return true
	}
	return false
}
