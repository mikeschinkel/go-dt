package dtx

import (
	"reflect"
)

func IsNilableKind(k reflect.Kind) (nilable bool) {
	switch k {
	case reflect.Ptr,
		reflect.Slice,
		reflect.Map,
		reflect.Chan,
		reflect.Func,
		reflect.Interface:
		nilable = true
	default:
		// Just here to stop Golang from complaining that switch is not exhaustive
	}
	return nilable
}

func IsNilableType(t reflect.Type) bool {
	return IsNilableKind(t.Kind())
}

func IsNilable(value any) (nilable bool) {
	if value == nil {
		goto end
	}
	nilable = IsNilableKind(reflect.ValueOf(value).Kind())
end:
	return
}

func IsNil(value any) (isNil bool) {
	var v reflect.Value
	if value == nil {
		isNil = true
		goto end
	}
	v = reflect.ValueOf(value)
	isNil = IsNilableKind(v.Kind()) && v.IsNil()
end:
	return isNil
}
