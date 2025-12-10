package dt

type URLSegment string

func (ps URLSegment) Contains(part any) bool {
	return EntryPath(ps).Contains(part)
}

func (ps URLSegment) HasDotDotPrefix() bool {
	return EntryPath(ps).HasDotDotPrefix()
}
