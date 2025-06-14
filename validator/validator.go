package validator

import (
	"fmt"
	"github.com/evenyosua18/ego/code"
	"github.com/go-playground/validator/v10"
	"sync"
)

var (
	val  *validator.Validate
	once sync.Once
)

func Validate(request any) error {
	// initiate
	once.Do(func() {
		val = validator.New()
	})

	// validate
	err := val.Struct(request)

	// manage error
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return code.Get(code.BadRequestError).
				SetErrorMessage(fmt.Sprintf("field '%s' failed on the tag is '%s'", e.Field(), e.ActualTag())).
				SetMessage("please check your request data again")
		}
	}

	return nil
}
