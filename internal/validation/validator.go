package validation

import (
	"github.com/anonychun/bibit/internal/api"
	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/gookit/validate"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewValidator)
}

type IValidator interface {
	Struct(s any) api.ValidationError
}

type Validator struct {
}

var _ IValidator = (*Validator)(nil)

func NewValidator(i do.Injector) (*Validator, error) {
	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false
		opt.SkipOnEmpty = false
		opt.FieldTag = "field"
	})

	return &Validator{}, nil
}

func (v *Validator) Struct(s any) api.ValidationError {
	validationErr := make(api.ValidationError)
	validation := validate.Struct(s)
	if !validation.Validate() {
		for field, messages := range validation.Errors {
			for _, message := range messages {
				validationErr.Add(field, message)
			}
		}
	}

	return validationErr
}
