package setup

import (
	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/go-xplat/env"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/os"
)

func New() *Setup {
	os := os.New()
	return &Setup{
		OS:      os,
		FS:      fs.NewOS(),
		Path:    filepath.NewProcessorWithOS(os),
		Env:     env.NewOS(),
		Console: console.NewOS(),
	}
}
