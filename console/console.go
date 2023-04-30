// Package console provides stdin, stderr, stdout, executable and args abstraction through an interface
package console

import (
	"io"
	"os"
)

type Console interface {
	In() io.Reader
	Out() io.Writer
	Error() io.Writer
	Executable() (string, error)
	Args() []string
}

type console struct {
}

func NewOS() Console {
	return &console{}
}

func (c *console) In() io.Reader {
	return os.Stdin
}

func (c *console) Out() io.Writer {
	return os.Stdout
}

func (c *console) Error() io.Writer {
	return os.Stderr
}

func (c *console) Args() []string {
	return os.Args
}

func (c *console) Executable() (string, error) {
	return os.Executable()
}
