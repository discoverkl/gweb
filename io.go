package gweb

import (
	"fmt"
	"io"
	"os"
	"syscall/js"
)

const StdinElementID = "stdin" // html input element

var Stdin io.Reader

func init() {
	input := js.Global().Get("document").Call("getElementById", StdinElementID)
	stdin, _ := NewReader(input) // stdin is optional, so we ignore error here
	Stdin = io.TeeReader(stdin, &echoWriter{w: os.Stdout})
}

func EchoStdin(prefix string) {
	echoInput = true
	echoPrefix = prefix
}

func NewReader(input js.Value) (io.Reader, error) {
	r := &reader{
		in: make(chan []byte, 8),
	}
	if input == js.Null() {
		return r, fmt.Errorf("can't find input element")
	}
	input.Call("addEventListener", "keydown", js.NewCallback(func(args []js.Value) {
		evt := args[0]
		code := evt.Get("keyCode").Int()
		if code == 13 {
			line := input.Get("value").String() + "\n"
			go func() {
				r.in <- []byte(line)
			}()
			input.Set("value", "")
			evt.Call("preventDefault")
		}
	}))
	return r, nil
}

type reader struct {
	pending []byte
	in      chan []byte // never close
}

func (r *reader) Read(p []byte) (n int, err error) {
	if len(r.pending) == 0 {
		r.pending = <-r.in
	}
	n = copy(p, r.pending)
	r.pending = r.pending[n:]
	return n, nil
}

var echoInput bool
var echoPrefix = "> "

type echoWriter struct {
	w     io.Writer
	dirty bool
}

// Write prefix each line with p.prefix.
// It assumes that a newline charactor always appears at the end of the input b.
// It respects the package level switch echoInput.
func (p *echoWriter) Write(b []byte) (int, error) {
	if !echoInput {
		return len(b), nil
	}

	if echoPrefix == "" {
		return p.w.Write(b)
	}

	if !p.dirty {
		p.w.Write([]byte(echoPrefix))
		p.dirty = true
	}
	n, err := p.w.Write(b)
	if n > 0 {
		if b[n-1] == '\n' {
			p.dirty = false
		}
	}
	return n, err
}
