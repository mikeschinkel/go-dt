package dt

type URLSegment string

func ParseURLSegment(s string) (uss URLSegment, err error) {
	if s == "" {
		err = NewErr(
			ErrInvalidURLSegment,
			ErrEmpty,
		)
		goto end
	}
	uss = URLSegment(s)
end:
	return uss, err
}

func (ps URLSegment) Contains(part any) bool {
	return EntryPath(ps).Contains(part)
}

func (ps URLSegment) HasDotDotPrefix() bool {
	return EntryPath(ps).HasDotDotPrefix()
}
