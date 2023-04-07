package fs

import (
	iofs "io/fs"
	"os"
)

type osfs struct {
}

func NewOS() FS {
	return &osfs{}
}

// Open implements FS
func (*osfs) Open(name string) (iofs.File, error) {
	return os.Open(name)
}

// Rename implements FS
func (*osfs) Rename(oldPath string, newPath string) error {
	return os.Rename(oldPath, newPath)
}

// Remove implements FS
func (*osfs) Remove(name string) error {
	return os.Remove(name)
}

// RemoveAll implements FS
func (*osfs) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// Glob implements FS
func (o *osfs) Glob(pattern string) ([]string, error) {
	return iofs.Glob(o, pattern)
}

// ReadDir implements FS
func (*osfs) ReadDir(name string) ([]iofs.DirEntry, error) {
	return os.ReadDir(name)
}

// ReadFile implements FS
func (*osfs) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

// Stat implements FS
func (*osfs) Stat(name string) (iofs.FileInfo, error) {
	return os.Stat(name)
}

// Sub implements FS
func (o *osfs) Sub(dir string) (iofs.FS, error) {
	return iofs.Sub(o, dir)
}
