package dtx

import (
	"fmt"
	"testing"

	"github.com/mikeschinkel/go-dt"
)

func TempTestDir(t *testing.T) dt.DirPath {
	t.Helper()
	return dt.DirPath(t.TempDir())
}

func SetTestEnv(t *testing.T, name string, value fmt.Stringer) {
	t.Helper()
	t.Setenv(name, value.String())
}
