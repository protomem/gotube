package validation

type Validator struct {
	Errors      []string          `json:"validationErrors,omitempty"`
	FieldErrors map[string]string `json:"validationFieldErrors,omitempty"`
}

func New() *Validator {
	return &Validator{}
}

func (v Validator) HasErrors() bool {
	return len(v.Errors) != 0 || len(v.FieldErrors) != 0
}

func (v *Validator) AddError(message string) {
	if v.Errors == nil {
		v.Errors = []string{}
	}

	v.Errors = append(v.Errors, message)
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = map[string]string{}
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) Check(ok bool, message string) {
	if !ok {
		v.AddError(message)
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func (v *Validator) Error() string {
	return "validation error(s)"
}

func (v *Validator) As(target any) bool {
	if _, ok := target.(*Validator); ok {
		return true
	}
	return false
}

func Validate(fn func(v *Validator)) error {
	v := New()
	fn(v)
	if v.HasErrors() {
		return v
	}
	return nil
}
