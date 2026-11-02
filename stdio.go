package dt

import (
	"io"
	"os"
)

var _ StdIO = (*stdIO)(nil)

type stdIO struct {
	out io.Writer
	err io.Writer
	in  io.Reader
}

func (s stdIO) Output() io.Writer {
	return s.out
}

func (s stdIO) ErrOut() io.Writer {
	return s.err
}

func (s stdIO) Input() io.Reader {
	return s.in
}

type StdIO interface {
	Output() io.Writer
	ErrOut() io.Writer
	Input() io.Reader
}

func NewDefaultStdIO() StdIO {
	return NewStdIO(os.Stdout, os.Stderr, os.Stdin)
}

func NewStdIO(out io.Writer, err io.Writer, in ...io.Reader) StdIO {
	var inp io.Reader
	switch len(in) {
	case 0:
		inp = nil
	default:
		inp = in[0]
	}
	return stdIO{
		out: out,
		err: err,
		in:  inp,
	}
}
