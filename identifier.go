package dt

type Identifier string

func ParseIdentifier(s string) (id Identifier, err error) {
	if s == "" {
		err = NewErr(
			ErrInvalidIdentifier,
			ErrEmpty,
		)
		goto end
	}
	id = Identifier(s)
end:
	return id, err
}
