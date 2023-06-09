package fs

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
	fstest "testing/fstest"

	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/platform"
)

type memory struct {
	fs        fstest.MapFS
	processor filepath.Processor
}

func NewMemory(options ...MemoryOption) FS {
	m := &memory{}
	for _, op := range options {
		op(m)
	}
	if m.processor == nil {
		m.processor = filepath.NewProcessor()
	}
	if m.fs == nil {
		m.fs = fstest.MapFS{}
	}
	return m
}

type MemoryOption = func(*memory)

func WithProcessor(processor filepath.Processor) MemoryOption {
	return func(m *memory) {
		m.processor = processor
	}
}

func WithPlatform(plat platform.Platform) MemoryOption {
	return func(m *memory) {
		m.processor = filepath.NewProcessorWithPlatform(plat)
		m.fs = fstest.MapFS{}
	}
}

func (m *memory) Create(name string) (File, error) {
	original := name
	name = m.normalizePath(name)

	file, ok := m.fs[name]
	if !ok {
		file = &fstest.MapFile{}
		m.fs[name] = file
	}
	file.Data = nil
	file.Mode = 0666
	return &openFile{
		path: original,
		infoFile: infoFile{
			name: m.processor.Base(original),
			file: file,
		},
	}, nil
}

func (m *memory) normalizePath(name string) string {
	if m.processor.Comparison() == filepath.IgnoreCase {
		return strings.ToLower(name)
	}
	return name
}

// Open implements FS
func (m *memory) Open(name string) (fs.File, error) {
	op := "open"
	original := name
	name = m.normalizePath(name)

	f, ok := m.fs[name]
	if !ok {
		return nil, &fs.PathError{
			Op:   op,
			Path: original,
			Err:  fs.ErrNotExist,
		}
	}
	return &openFile{
		path: name,
		infoFile: infoFile{
			name: m.processor.Base(name),
			file: f,
		},
	}, nil
}

func isReadOnly(mode int) bool {
	switch {
	case mode&os.O_APPEND == os.O_APPEND:
		return false
	case mode&os.O_CREATE == os.O_CREATE:
		return false
	case mode&os.O_TRUNC == os.O_CREATE:
		return false
	case mode&os.O_WRONLY == os.O_WRONLY:
		return false
	case mode&os.O_RDWR == os.O_RDWR:
		return false
	}
	return true
}

// OpenFile implements OpenFS
func (m *memory) OpenFile(name string, mode int, perm fs.FileMode) (File, error) {
	op := "openFile"
	original := name
	name = m.normalizePath(name)

	f, ok := m.fs[name]
	if !ok {
		// for readonly files, if the file doesn't exist return an error
		if isReadOnly(mode) {
			return nil, &fs.PathError{
				Op:   op,
				Path: original,
				Err:  fs.ErrNotExist,
			}
		}

		f = &fstest.MapFile{}
		m.fs[name] = f
	}

	// truncate if O_TRUNC specified
	if mode&os.O_TRUNC != 0 {
		f.Data = nil
	}

	// seek pos
	offset := 0
	if mode&os.O_APPEND != 0 {
		offset = len(f.Data)
	}

	return &openFile{
		path:   name,
		offset: int64(offset),
		infoFile: infoFile{
			name: m.processor.Base(name),
			file: f,
		},
	}, nil
}

// Rename implements FS
func (m *memory) Rename(oldPath string, newPath string) error {

	oldPath = m.normalizePath(oldPath)
	newPath = m.normalizePath(newPath)

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
	path = m.normalizePath(path)
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
	for p := range m.fs {
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
	// check that we can open the file
	d, err := m.Open(name)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	// create the list of entries
	var entries []fs.DirEntry
	for path, file := range m.fs {
		originalPath := path
		// make sure both are lowered if case insensitive
		if m.processor.Comparison() == filepath.IgnoreCase {
			path = strings.ToLower(path)
			name = strings.ToLower(name)
		}

		// same dir
		if path == name {
			continue
		}

		// any file will have the prefix of the path
		if strings.HasPrefix(path, name) {

			// get the file name
			fileName := m.processor.Base(originalPath)

			// append
			entries = append(entries, &infoFile{name: fileName, file: file})
		}
	}
	return entries, nil
}

// ReadFile implements FS
func (m *memory) ReadFile(name string) ([]byte, error) {

	f, err := m.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, stat.Size())

	_, err = f.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// WriteFile implements FS
func (m *memory) WriteFile(name string, data []byte, perm os.FileMode) error {
	name = m.normalizePath(name)
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
	path = m.normalizePath(path)
	_, ok := m.fs[path]
	return ok, nil
}

// Stat implements FS
func (m *memory) Stat(name string) (fs.FileInfo, error) {
	f, err := m.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.Stat()
}

// Sub implements FS
func (m *memory) Sub(dir string) (fs.FS, error) {
	return m.fs.Sub(dir)
}

// Mkdir implements MakeDirFS
func (m *memory) Mkdir(path string, perm fs.FileMode) error {
	fp, err := m.processor.Parser().Parse(path)
	if err != nil {
		return err
	}
	accumulator := fp.Root()

	// check each ancestor path
	for i := 0; i < len(fp.Segments); i++ {
		currentPath := accumulator.String(m.processor.Separator())
		currentPath = m.normalizePath(currentPath)
		_, ok := m.fs[currentPath]
		if !ok {
			return errNotExist(currentPath)
		}
		seg := fp.Segments[i]
		fpseg, err := m.processor.Parser().Parse(seg)
		if err != nil {
			return err
		}
		accumulator = accumulator.Join(fpseg)
	}

	// write the segment
	m.fs[path] = &fstest.MapFile{
		Mode: perm | fs.ModeDir,
	}

	return nil
}

// MkdirAll implements MakeDirFS
func (m *memory) MkdirAll(path string, perm fs.FileMode) error {
	fp, err := m.processor.Parser().Parse(path)
	if err != nil {
		return err
	}
	accumulator := fp.Root()

	// create each ancestor path
	for i := 0; i < len(fp.Segments); i++ {
		currentPath := accumulator.String(m.processor.Separator())
		currentPath = m.normalizePath(currentPath)
		_, ok := m.fs[currentPath]
		if !ok {
			m.fs[currentPath] = &fstest.MapFile{
				Mode: perm | fs.ModeDir,
			}
		}
		seg := fp.Segments[i]
		fpseg, err := m.processor.Parser().Parse(seg)
		if err != nil {
			return err
		}
		accumulator = accumulator.Join(fpseg)
	}

	// create the path
	m.fs[path] = &fstest.MapFile{
		Mode: perm | fs.ModeDir,
	}

	return nil
}

func errNotExist(path string) error {
	return fmt.Errorf("'%s' %w", path, fs.ErrNotExist)
}
