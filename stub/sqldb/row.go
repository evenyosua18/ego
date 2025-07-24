package sqldb

import (
	"fmt"
	"reflect"
)

type StubSqlRow struct {
	Values []any
	Err    error
}

func (s *StubSqlRow) Scan(dest ...any) error {
	if s.Err != nil {
		return s.Err
	}

	if len(dest) != len(s.Values) {
		return fmt.Errorf("scan mismatch: got %d destinations, but %d values", len(dest), len(s.Values))
	}

	for i, val := range s.Values {
		// Use reflection to set the value if it's a pointer
		ptrVal := dest[i]
		if ptr, ok := ptrVal.(*any); ok {
			*ptr = val
		} else {
			// type assert to pointer and set via reflection
			rv := reflect.ValueOf(ptrVal)
			if rv.Kind() != reflect.Ptr {
				return fmt.Errorf("destination at index %d is not a pointer", i)
			}
			rv.Elem().Set(reflect.ValueOf(val))
		}
	}

	return nil
}
