package dtx

import (
	"github.com/mikeschinkel/go-dt"
)

func EntryStatusError(status dt.EntryStatus) (err error) {
	switch status {
	case dt.IsMissingEntry:
		err = NewErr(dt.ErrFileNotExists)
	case dt.IsDirEntry:
		err = NewErr(dt.ErrPathIsDir)
	case dt.IsSocketEntry, dt.IsPipeEntry, dt.IsDeviceEntry:
		err = NewErr(dt.ErrUnsupportedEntryType)
	case dt.IsInvalidEntryStatus:
		err = NewErr(dt.ErrInvalidEntryStatus)
	case dt.IsUnclassifiedEntryStatus:
		fallthrough
	default:
		err = NewErr(dt.ErrUnclassifiedEntryStatus)
	}
	return err
}
