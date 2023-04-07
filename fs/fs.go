package fs

import (
	iofs "io/fs"
)

type RenameFS interface {
	iofs.FS
	Rename(oldPath, newPath string) error
}

type RemoveFS interface {
	iofs.FS
	Remove(name string) error
}

type FS interface {
	iofs.FS
	RenameFS
	RemoveFS
	iofs.GlobFS
	iofs.ReadDirFS
	iofs.ReadFileFS
	iofs.ReadFileFS
	iofs.StatFS
	iofs.SubFS
}
