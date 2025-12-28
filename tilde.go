package dt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func DirPathToTilde(path DirPath, option TildeOption) (tdp TildeDirPath) {
	return ToTilde[DirPath, TildeDirPath](path, option)
}

func FilepathToTilde(path Filepath, option TildeOption) (tdp TildeFilepath) {
	return ToTilde[Filepath, TildeFilepath](path, option)
}

func EntryPathToTilde(path EntryPath, option TildeOption) (tdp TildeEntryPath) {
	return ToTilde[EntryPath, TildeEntryPath](path, option)
}

type tildable interface {
	DirPath | Filepath | EntryPath
}

type tilde interface {
	TildeDirPath | TildeFilepath | TildeEntryPath
}

type TildeOption uint8

func (t TildeOption) String() string {
	switch t {
	case OrEmptyString:
		return "OrEmptyString"
	case OrFullPath:
		return "OrFullPath"
	case OrPanic:
		return "OrPanic"
	case UnspecifiedTildeOption:
		return "Unspecified"
	default:
		return "Invalid"
	}
}

const (
	UnspecifiedTildeOption TildeOption = 0
	OrEmptyString          TildeOption = 1
	OrFullPath             TildeOption = 2
	OrPanic                TildeOption = 3
)

func ToTilde[P tildable, TP tilde](path P, option TildeOption) (tp TP) {
	var rel string
	var err error
	var homeDir DirPath

	switch option {
	case OrEmptyString, OrFullPath, OrPanic:
		// Continue on
	case UnspecifiedTildeOption:
		fallthrough
	default:
		panic(fmt.Sprintf("dt.ToTilde() called with invalid TildeOption; must be %s, %s, or %s",
			OrEmptyString,
			OrFullPath,
			OrPanic,
		))
	}

	if len(path) == 0 {
		goto end
	}
	switch {
	case path == "~":
		fallthrough
	case len(path) >= 2 && string(path[:2]) == "~"+string(os.PathSeparator):
		tp = TP(path)
		goto end
	}

	homeDir = GetUserHomeDir()
	rel, err = filepath.Rel(string(homeDir), string(path))

	switch {
	case err != nil:
		fallthrough
	case rel == "..":
		fallthrough
	case strings.HasPrefix(rel, ".."+string(os.PathSeparator)):
		switch option {
		case OrEmptyString:
			tp = ""
			goto end
		case OrFullPath:
			tp = TP(path)
			goto end
		default:
			panic(fmt.Sprintf("%s must be within HOME directory %s when calling dt.ToTilde()", path, homeDir))
		}
	}

	if rel == "." {
		tp = TP("~")
		goto end
	}

	tp = TP(fmt.Sprintf("~%c%s", os.PathSeparator, rel))
end:
	return tp
}

func TrimTilde[P tildable](path P) (ps PathSegments) {
	if len(path) < 2 {
		ps = PathSegments(path)
		goto end
	}
	if string(path[:2]) == "~"+string(os.PathSeparator) {
		ps = PathSegments(path[2:])
		goto end
	}
	ps = PathSegments(path)
end:
	return ps
}
