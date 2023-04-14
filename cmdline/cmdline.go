// Package cmdline provides stdin, stderr, stdout, executable and args abstraction through an interface
package cmdline

import (
	"io"
	"os"
)

type CommandLine interface {
	In() io.Reader
	Out() io.Writer
	Error() io.Writer
	Executable() (string, error)
	Args() []string
}

type commandLine struct {
}

func NewOS() CommandLine {
	return &commandLine{}
}

func (c *commandLine) In() io.Reader {
	return os.Stdin
}

func (c *commandLine) Out() io.Writer {
	return os.Stdout
}

func (c *commandLine) Error() io.Writer {
	return os.Stderr
}

func (c *commandLine) Args() []string {
	return os.Args
}

func (c *commandLine) Executable() (string, error) {
	return os.Executable()
}
