package sqldb

import (
	"fmt"
	"reflect"
)

type StubSqlRows struct {
	Values      []any
	Destination any

	CloseErr     error
	MapResultErr error
	ScanAllErr   error
	ScanOneErr   error

	RowIndex int
}

func (s *StubSqlRows) Close() error {
	return s.CloseErr
}

func (s *StubSqlRows) Next() bool {
	if s.RowIndex < len(s.Values) {
		s.RowIndex++
		return true
	}
	return false
}

func (s *StubSqlRows) MapResult(dest ...any) error {
	if s.MapResultErr != nil {
		return s.MapResultErr
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

func (s *StubSqlRows) ScanAll(dest any) error {
	if s.ScanAllErr != nil {
		return s.ScanAllErr
	}

	// dest should be a pointer to a slice
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr || destVal.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("ScanAll expects a pointer to a slice, got %T", dest)
	}

	sliceVal := destVal.Elem()
	for _, val := range s.Values {
		sliceVal.Set(reflect.Append(sliceVal, reflect.ValueOf(val)))
	}
	return nil
}

func (s *StubSqlRows) ScanOne(dest any) error {
	if s.ScanOneErr != nil {
		return s.ScanOneErr
	}

	// dest should be a pointer to struct or value
	if len(s.Values) == 0 {
		return fmt.Errorf("no data to scan")
	}

	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("ScanOne expects a pointer, got %T", dest)
	}
	val.Elem().Set(reflect.ValueOf(s.Values[0]))
	return nil
}
