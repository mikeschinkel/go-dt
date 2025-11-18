package dtx

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/mikeschinkel/go-dt"
	"github.com/mikeschinkel/go-dt/de"
)

var ErrAssertTypeFailed = errors.New("assert type failed")

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
		err = dt.NewErr(dt.ErrValueIsNil)
		goto end
	}

	t, ok = value.(T)
	if !ok {
		err = dt.NewErr(dt.ErrFailedTypeAssertion)
		goto end
	}

	// Enforce: if T is nilable, don't allow a typed-nil result on success.
	if IsNil(t) {
		err = dt.NewErr(dt.ErrInterfaceValueIsNil)
		goto end
	}

end:
	if err != nil {
		err = dt.WithErr(err,
			ErrAssertTypeFailed,
			"target", reflect.TypeOf((*T)(nil)).Elem().String(),
			"source", fmt.Sprintf("%T", value),
		)
	}
	return t, err
}
