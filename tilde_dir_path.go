package dt

type TildeDirPath string

func ParseTildeDirPath(s string) (tdp TildeDirPath, err error) {
	var tep TildeEntryPath
	tep, err = ParseTildeEntryPath(s)
	if err != nil {
		goto end
	}
	tdp = TildeDirPath(tep)
end:
	return tdp, err
}

func (tdp TildeDirPath) Expand() (_ DirPath, err error) {
	var ep EntryPath
	ep, err = EntryPath(tdp).Expand()
	return DirPath(ep), err
}
