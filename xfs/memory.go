package xfs

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
	fstest "testing/fstest"

	"github.com/patrickhuber/go-xplat/xfilepath"
)

type memory struct {
	fs            fstest.MapFS
	pathSeperator string
}

func NewMemory(options ...MemoryOption) FS {
	m := &memory{
		fs:            fstest.MapFS{},
		pathSeperator: "/",
	}
	for _, op := range options {
		op(m)
	}
	return m
}

type MemoryOption = func(*memory)

func WithPathSeperator(sep string) MemoryOption {
	return func(m *memory) {
		m.pathSeperator = sep
	}
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

// WriteFile implements FS
func (m *memory) WriteFile(name string, data []byte, perm os.FileMode) error {
	file, ok := m.fs[name]
	if !ok {
		file = &fstest.MapFile{}
		m.fs[name] = file
	}
	file.Data = data
	file.Mode = perm
	return nil
}

// Exists implements FS
func (m *memory) Exists(path string) (bool, error) {
	_, ok := m.fs[path]
	return ok, nil
}

// Stat implements FS
func (m *memory) Stat(name string) (fs.FileInfo, error) {
	return m.fs.Stat(name)
}

// Sub implements FS
func (m *memory) Sub(dir string) (fs.FS, error) {
	return m.fs.Sub(dir)
}

// Mkdir implements MakeDirFS
func (m *memory) Mkdir(path string, perm fs.FileMode) error {
	fp, err := xfilepath.Parse(path)
	if err != nil {
		return err
	}

	root := m.root(fp)
	if root == "" {
		return nil
	}

	ok, err := m.Exists(root)
	if err != nil {
		return err
	}
	if !ok {
		if len(fp.Segments) > 0 {
			return errNotExist(root)
		}
		m.fs[root] = &fstest.MapFile{
			Mode: perm | fs.ModeDir,
		}
		return nil
	}

	var builder strings.Builder
	builder.WriteString(root)
	for i := 0; i < len(fp.Segments); i++ {

		seg := fp.Segments[i]
		builder.WriteString(m.pathSeperator)
		builder.WriteString(seg)
		p := builder.String()

		isLastSegment := i == len(fp.Segments)-1

		if isLastSegment {
			m.fs[p] = &fstest.MapFile{
				Mode: perm | fs.ModeDir,
			}
			return nil
		} else {
			ok, err := m.Exists(p)
			if err != nil {
				return err
			}
			if !ok {
				return errNotExist(p)
			}
		}
	}
	return nil
}

// MkdirAll implements MakeDirFS
func (m *memory) MkdirAll(path string, perm fs.FileMode) error {
	fp, err := xfilepath.Parse(path)
	if err != nil {
		return err
	}

	root := m.root(fp)
	if root == "" {
		return nil
	}

	err = m.Mkdir(root, perm)
	if err != nil {
		return err
	}

	var builder strings.Builder
	builder.WriteString(root)
	for _, seg := range fp.Segments {
		builder.WriteString(m.pathSeperator)
		builder.WriteString(seg)
		m.fs[builder.String()] = &fstest.MapFile{
			Mode: perm | fs.ModeDir,
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *memory) root(fp xfilepath.FilePath) string {
	switch {
	case !fp.Absolute:
		// Relative Path
		return ""

	case fp.Volume.Drive != "":
		// windows path
		return fp.Volume.Drive

	case fp.Volume.Host != "":
		// UNC Path
		return volumeName(fp.Volume, m.pathSeperator)

	default:
		// unix path
		return string(m.pathSeperator)
	}

}

func errNotExist(path string) error {
	return fmt.Errorf("'%s' %w", path, fs.ErrNotExist)
}

func volumeName(v xfilepath.Volume, sep string) string {
	var builder strings.Builder
	builder.WriteString(sep)
	builder.WriteString(sep)
	builder.WriteString(v.Host)
	builder.WriteString(sep)
	builder.WriteString(v.Share)
	return builder.String()
}
