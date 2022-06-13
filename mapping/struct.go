package mapping

import (
	"errors"
	"reflect"
)

const (
	ValueError = `value must be struct`
)

func InterfaceToStruct(in interface{}) (out interface{}, err error) {
	if reflect.TypeOf(in).Kind() == reflect.Struct {

	} else {
		return nil, errors.New(ValueError)
	}
	return nil, nil
}
