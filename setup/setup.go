package setup

import (
	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/go-xplat/env"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/os"
)

type Setup struct {
	OS      os.OS
	FS      fs.FS
	Path    *filepath.Processor
	Env     env.Environment
	Console console.Console
}
