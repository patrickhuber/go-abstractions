package host

import (
	"github.com/patrickhuber/go-xplat/arch"
	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/go-xplat/env"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/patrickhuber/go-xplat/platform"
)

type Host struct {
	OS      os.OS
	FS      fs.FS
	Path    *filepath.Processor
	Env     env.Environment
	Console console.Console
}

func New() *Host {
	os := os.New()
	return &Host{
		OS:      os,
		FS:      fs.NewOS(),
		Path:    filepath.NewProcessorWithOS(os),
		Env:     env.NewOS(),
		Console: console.NewOS(),
	}
}

func NewTest(
	p platform.Platform,
	a arch.Arch,
	vars map[string]string,
	args []string) *Host {

	os := os.NewMock(
		os.WithArchitecture(a),
		os.WithPlatform(p))

	path := filepath.NewProcessorWithOS(os)
	return &Host{
		OS:      os,
		Path:    path,
		FS:      fs.NewMemory(fs.WithProcessor(path)),
		Env:     env.NewMemoryWithMap(vars),
		Console: console.NewMemory(console.WithArgs(args)),
	}
}