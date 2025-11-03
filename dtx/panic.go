package dtx

import (
	"fmt"
)

func Panicf(format string, args ...any) {
	panic(fmt.Sprintf(format, args...))
}
