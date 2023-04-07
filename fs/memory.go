package fs

import (
	"io/fs"
	"os"
	"strings"
	fstest "testing/fstest"
)

type memory struct {
	fs fstest.MapFS
}

// Open implements FS
func (m *memory) Open(name string) (fs.File, error) {
	return m.fs.Open(name)
}

// Rename implements FS
func (m *memory) Rename(oldPath string, newPath string) error {
	file, ok := m.fs[oldPath]
	if !ok {
		return os.ErrNotExist
	}
	delete(m.fs, oldPath)
	m.fs[newPath] = file
	return nil
}

// Remove implements FS
func (m *memory) Remove(path string) error {
	_, ok := m.fs[path]
	if !ok {
		return os.ErrNotExist
	}
	delete(m.fs, path)
	return nil
}

// RemoveAll implements FS
func (m *memory) RemoveAll(path string) error {
	paths := []string{}
	for p, _ := range m.fs {
		if strings.HasPrefix(p, path) {
			paths = append(paths, p)
		}
	}
	for _, p := range paths {
		delete(m.fs, p)
	}
	return nil
}

// Glob implements FS
func (m *memory) Glob(pattern string) ([]string, error) {
	return m.fs.Glob(pattern)
}

// ReadDir implements FS
func (m *memory) ReadDir(name string) ([]fs.DirEntry, error) {
	return m.fs.ReadDir(name)
}

// ReadFile implements FS
func (m *memory) ReadFile(name string) ([]byte, error) {
	return m.fs.ReadFile(name)
}

// Stat implements FS
func (m *memory) Stat(name string) (fs.FileInfo, error) {
	return m.fs.Stat(name)
}

// Sub implements FS
func (m *memory) Sub(dir string) (fs.FS, error) {
	return m.fs.Sub(dir)
}

func NewMemory() FS {
	return &memory{
		fs: fstest.MapFS{},
	}
}
