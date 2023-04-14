package standard

import (
	"io"
	"os"
)

type Standard interface {
	In() io.Reader
	Out() io.Writer
	Error() io.Writer
}

type osStandard struct {
}

func NewOS() Standard {
	return &osStandard{}
}

func (c *osStandard) In() io.Reader {
	return os.Stdin
}

func (c *osStandard) Out() io.Writer {
	return os.Stdout
}

func (c *osStandard) Error() io.Writer {
	return os.Stderr
}
