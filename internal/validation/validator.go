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

type Validator interface {
	Struct(s any) api.ValidationError
}

type ValidatorImpl struct {
}

func NewValidator(i do.Injector) (Validator, error) {
	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false
		opt.SkipOnEmpty = false
		opt.FieldTag = "field"
	})

	return &ValidatorImpl{}, nil
}

func (v *ValidatorImpl) Struct(s any) api.ValidationError {
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
