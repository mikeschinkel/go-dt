package de

import (
	"errors"
)

var (
	ErrFailedTypeAssertion     = errors.New("failed type assertion")
	ErrFailedReadingSymlink    = errors.New("failed reading symlink")
	ErrFailedToLoadFile        = errors.New("failed to load file")
	ErrFailedCreatingDirectory = errors.New("failed to create directory")
	ErrContainsBackslash       = errors.New("contains slash ('\\')")
	ErrContainsSlash           = errors.New("contains slash ('/')")
	ErrEmpty                   = errors.New("empty")
	ErrInvalidPercentEncoding  = errors.New("invalid percent encoding")
	ErrTooLong                 = errors.New("too long")
	ErrUnspecified             = errors.New("unspecified")
	ErrInvalid                 = errors.New("invalid")
	ErrControlCharacter        = errors.New("control character")
	ErrInvalidCharacter        = errors.New("invalid charnacter")
	ErrTrailingSpace           = errors.New("trailing space")
	ErrTrailingPeriod          = errors.New("trailing period")
	ErrReservedDeviceName      = errors.New("reserved device name")
)
