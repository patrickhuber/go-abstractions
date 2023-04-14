package standard

import (
	"bytes"
	"io"
)

type memory struct {
	in  *bytes.Buffer
	out *bytes.Buffer
	err *bytes.Buffer
}

type MemoryStandard interface {
	Standard
	// OutBuffer exposes the output buffer for the memory console to enable testing
	OutBuffer() *bytes.Buffer
	// ErrBuffer exposes the error buffer for the memory console to enable testing
	ErrBuffer() *bytes.Buffer
	// InBuffer exposes the error buffer for the memory console to enable testing
	InBuffer() *bytes.Buffer
}

func NewMemory() MemoryStandard {
	return &memory{
		in:  &bytes.Buffer{},
		out: &bytes.Buffer{},
		err: &bytes.Buffer{},
	}
}

func (c *memory) In() io.Reader {
	return c.in
}

func (c *memory) Out() io.Writer {
	return c.out
}

func (c *memory) Error() io.Writer {
	return c.err
}

// ErrBuffer implements MemoryConsole
func (c *memory) ErrBuffer() *bytes.Buffer {
	return c.err
}

// InBuffer implements MemoryConsole
func (c *memory) InBuffer() *bytes.Buffer {
	return c.in
}

// OutBuffer implements MemoryConsole
func (c *memory) OutBuffer() *bytes.Buffer {
	return c.out
}
