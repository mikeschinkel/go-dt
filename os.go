package dt

import (
	"io/fs"
	"os"
)

func DirFS(dp DirPath) fs.FS {
	return os.DirFS(string(dp))
}

func MkdirAll(dp DirPath, mode os.FileMode) error {
	return os.MkdirAll(string(dp), mode)
}

func MkdirTemp(dp DirPath, pattern string) (DirPath, error) {
	td, err := os.MkdirTemp(string(dp), pattern)
	return DirPath(td), err
}

func RemoveAll(dp DirPath) error {
	return os.RemoveAll(string(dp))
}

func UserHomeDir() (DirPath, error) {
	dp, err := os.UserHomeDir()
	return DirPath(dp), err
}

func CreateFile(fp Filepath) (*os.File, error) {
	return os.Create(string(fp))
}

func ReadFile(fp Filepath) ([]byte, error) {
	return os.ReadFile(string(fp))
}

func WriteFile(fp Filepath, data []byte, mode os.FileMode) error {
	return os.WriteFile(string(fp), data, mode)
}

func UserConfigDir() (DirPath, error) {
	cd, err := os.UserConfigDir()
	return DirPath(cd), err
}

func Getwd() (DirPath, error) {
	wd, err := os.Getwd()
	return DirPath(wd), err
}
