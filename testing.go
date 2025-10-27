package dt

import (
	"testing"
)

func TempDir(t *testing.T) DirPath {
	t.Helper()
	return DirPath(t.TempDir())
}
