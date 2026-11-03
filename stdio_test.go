package dt

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestNewStdIO(t *testing.T) {
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	in := strings.NewReader("input")

	tests := []struct {
		name    string
		in      []io.Reader
		wantIn  io.Reader
		wantNil bool
	}{
		{name: "with input", in: []io.Reader{in}, wantIn: in},
		{name: "without input", in: nil, wantNil: true},
		{name: "first input wins", in: []io.Reader{in, strings.NewReader("other")}, wantIn: in},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStdIO(out, errOut, tt.in...)
			if s.Output() != io.Writer(out) {
				t.Fatalf("Output() = %v, want %v", s.Output(), out)
			}
			if s.ErrOut() != io.Writer(errOut) {
				t.Fatalf("ErrOut() = %v, want %v", s.ErrOut(), errOut)
			}
			if tt.wantNil {
				if s.Input() != nil {
					t.Fatalf("Input() = %v, want nil", s.Input())
				}
				return
			}
			if s.Input() != tt.wantIn {
				t.Fatalf("Input() = %v, want %v", s.Input(), tt.wantIn)
			}
		})
	}
}

func TestNewStdIOStreamsAreUsable(t *testing.T) {
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	s := NewStdIO(out, errOut, strings.NewReader("hello"))

	if _, err := io.WriteString(s.Output(), "out"); err != nil {
		t.Fatalf("Output() write error = %v", err)
	}
	if _, err := io.WriteString(s.ErrOut(), "err"); err != nil {
		t.Fatalf("ErrOut() write error = %v", err)
	}
	got, err := io.ReadAll(s.Input())
	if err != nil {
		t.Fatalf("Input() read error = %v", err)
	}

	if out.String() != "out" {
		t.Fatalf("Output() captured %q, want %q", out.String(), "out")
	}
	if errOut.String() != "err" {
		t.Fatalf("ErrOut() captured %q, want %q", errOut.String(), "err")
	}
	if string(got) != "hello" {
		t.Fatalf("Input() read %q, want %q", got, "hello")
	}
}

func TestNewDefaultStdIO(t *testing.T) {
	s := NewDefaultStdIO()
	if s.Output() != io.Writer(os.Stdout) {
		t.Fatalf("Output() is not os.Stdout")
	}
	if s.ErrOut() != io.Writer(os.Stderr) {
		t.Fatalf("ErrOut() is not os.Stderr")
	}
	if s.Input() != io.Reader(os.Stdin) {
		t.Fatalf("Input() is not os.Stdin")
	}
}
