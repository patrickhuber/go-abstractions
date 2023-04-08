package fs

import (
	iofs "io/fs"
	"os"
)

type RenameFS interface {
	iofs.FS
	Rename(oldPath, newPath string) error
}

type RemoveFS interface {
	iofs.FS
	Remove(name string) error
}

type WriteFileFS interface {
	WriteFile(name string, data []byte, perm os.FileMode) error
}

type ExistsFS interface {
	Exists(path string) (bool, error)
}

type FS interface {
	iofs.FS
	RenameFS
	RemoveFS
	WriteFileFS
	ExistsFS
	iofs.GlobFS
	iofs.ReadDirFS
	iofs.ReadFileFS
	iofs.ReadFileFS
	iofs.StatFS
	iofs.SubFS
}
