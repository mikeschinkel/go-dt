package dtx

import (
	"fmt"
	"reflect"

	"github.com/mikeschinkel/go-dt"
	"github.com/mikeschinkel/go-dt/de"
)

// AssertType performs a safe type assertion from any to T.
// Returns a typed value and a dt.NewErr() error when the assertion fails.
//
// Typical usage:
//
//	ro, err := dtx.AssertType[*RawOptions](opts)
//
// If value is nil or of a non-matching type, err is non-nil and t is zero.
func AssertType[T any](value any) (t T, err error) {
	var ok bool

	if value == nil {
		err = dt.NewErr(
			de.ErrFailedTypeAssertion,
			"source", "nil",
		)
		goto end
	}

	t, ok = value.(T)
	if !ok {
		err = dt.NewErr(
			de.ErrFailedTypeAssertion,
			"source", fmt.Sprintf("%T", value),
		)
	}

end:
	if err != nil {
		err = dt.WithErr(err,
			"target", reflect.TypeOf((*T)(nil)).Elem().String(),
		)
	}
	return t, err
}
