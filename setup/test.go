package setup

import (
	"github.com/patrickhuber/go-xplat/arch"
	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/go-xplat/env"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/patrickhuber/go-xplat/platform"
)

type Option func(*Options)

type Options struct {
	args     []string
	vars     map[string]string
	platform platform.Platform
	arch     arch.Arch
}

func Args(arg ...string) Option {
	return func(o *Options) {
		o.args = arg
	}
}

func Vars(vars map[string]string) Option {
	return func(o *Options) {
		o.vars = vars
	}
}

func Platform(p platform.Platform) Option {
	return func(o *Options) {
		o.platform = p
	}
}

func Arch(a arch.Arch) Option {
	return func(o *Options) {
		o.arch = a
	}
}

func NewTest(options ...Option) *Setup {
	op := &Options{}
	for _, option := range options {
		option(op)
	}
	if string(op.platform) == "" {
		op.platform = platform.Default()
	}
	if string(op.arch) == "" {
		op.arch = arch.AMD64
	}
	if op.vars == nil {
		op.vars = map[string]string{}
	}
	if op.args == nil {
		op.args = []string{}
	}
	os := os.NewMock(
		os.WithArchitecture(op.arch),
		os.WithPlatform(op.platform))
	path := filepath.NewProcessorWithOS(os)
	return &Setup{
		OS:      os,
		FS:      fs.NewMemory(fs.WithProcessor(path)),
		Path:    path,
		Env:     env.NewMemoryWithMap(op.vars),
		Console: console.NewMemory(console.WithArgs(op.args)),
	}
}
