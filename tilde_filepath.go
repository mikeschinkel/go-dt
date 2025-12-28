package dt

type TildeFilepath string

func ParseTildeFilepath(s string) (tdp TildeFilepath, err error) {
	var tep TildeEntryPath
	tep, err = ParseTildeEntryPath(s)
	if err != nil {
		goto end
	}
	tdp = TildeFilepath(tep)
end:
	return tdp, err
}

func (tdp TildeFilepath) Expand() (_ Filepath, err error) {
	var ep EntryPath
	ep, err = EntryPath(tdp).Expand()
	return Filepath(ep), err
}
