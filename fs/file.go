package fs

import (
	"io"
	"io/fs"
	"testing/fstest"
	"time"
)

type infoFile struct {
	name string
	file *fstest.MapFile
}

func (i *infoFile) Name() string               { return i.name }
func (i *infoFile) Size() int64                { return int64(len(i.file.Data)) }
func (i *infoFile) Mode() fs.FileMode          { return i.file.Mode }
func (i *infoFile) Type() fs.FileMode          { return i.file.Mode.Type() }
func (i *infoFile) ModTime() time.Time         { return i.file.ModTime }
func (i *infoFile) IsDir() bool                { return i.file.Mode&fs.ModeDir != 0 }
func (i *infoFile) Sys() any                   { return i.file.Sys }
func (i *infoFile) Info() (fs.FileInfo, error) { return i, nil }

type openFile struct {
	path string
	infoFile
	offset int64
}

func (f *openFile) Stat() (fs.FileInfo, error) {
	return f, nil
}

func (f *openFile) Close() error {
	return nil
}

func (f *openFile) Read(b []byte) (int, error) {
	op := "read"
	if f.file.Mode&fs.ModeDir != 0 {
		return 0, &fs.PathError{Op: op, Path: f.path, Err: fs.ErrInvalid}
	}
	if f.offset >= int64(len(f.file.Data)) {
		return 0, io.EOF
	}
	if f.offset < 0 {
		return 0, &fs.PathError{Op: op, Path: f.path, Err: fs.ErrInvalid}
	}
	n := copy(b, f.file.Data[f.offset:])
	f.offset += int64(n)
	return n, nil
}

func (f *openFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case 0:
		// offset += 0
	case 1:
		offset += f.offset
	case 2:
		offset += int64(len(f.file.Data))
	}
	if offset < 0 || offset > int64(len(f.file.Data)) {
		return 0, &fs.PathError{Op: "seek", Path: f.path, Err: fs.ErrInvalid}
	}
	f.offset = offset
	return offset, nil
}

func (f *openFile) ReadAt(b []byte, offset int64) (int, error) {
	if offset < 0 || offset > int64(len(f.file.Data)) {
		return 0, &fs.PathError{Op: "read", Path: f.path, Err: fs.ErrInvalid}
	}
	n := copy(b, f.file.Data[offset:])
	if n < len(b) {
		return n, io.EOF
	}
	return n, nil
}
