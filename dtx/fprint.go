package dtx

import (
	"fmt"
	"io"

	"github.com/mikeschinkel/go-dt"
)

func Fprintf(w io.Writer, format string, a ...any) (n int) {
	n, err := fmt.Fprintf(w, format, a...)
	if err != nil {
		dt.Logf("Failed to Fprintf: %s", err)
	}
	return n
}

func Fprintln(w io.Writer, a ...any) (n int) {
	n, err := fmt.Fprintln(w, a...)
	if err != nil {
		dt.Logf("Failed to Fprintln: %s", err)
	}
	return n
}
