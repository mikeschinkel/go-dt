package dtx

import (
	"reflect"
)

func IsZero(value any) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	// Automatically dereference pointers
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	return v.IsZero()
}
