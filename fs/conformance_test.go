package fs_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/stretchr/testify/require"
)

type conformance struct {
	fs   fs.FS
	path *filepath.Processor
}

type file struct {
	name    string
	content []byte
}

func NewConformance(fs fs.FS) *conformance {
	return &conformance{
		fs: fs,
	}
}

func NewConformanceWithPath(fs fs.FS, path *filepath.Processor) *conformance {
	return &conformance{
		fs:   fs,
		path: path,
	}
}

func (c *conformance) TestMkdirCreatesRoot(t *testing.T, root string) {
	err := c.fs.Mkdir(root, 0666)
	require.Nil(t, err)
	ok, err := c.fs.Exists(root)
	require.Nil(t, err)
	require.True(t, ok)
}

func (c *conformance) TestMkdirFailsWhenRootNotExists(t *testing.T, path string) {
	err := c.fs.Mkdir(path, 0666)
	require.NotNil(t, err)
}

func (c *conformance) TestMkdirAllCreatesAllDirectories(t *testing.T, path string, expected []string) {
	err := c.fs.MkdirAll(path, 0666)
	require.Nil(t, err)
	for _, p := range expected {
		ok, err := c.fs.Exists(p)
		require.Nil(t, err)
		require.True(t, ok, "%s does not exist", p)
	}
}

func (c *conformance) TestWriteFile(t *testing.T, path string, name string, content string) {

	require.NotNil(t, c.path)

	filep := c.path.Join(path, name)
	err := c.fs.MkdirAll(path, 0666)
	require.Nil(t, err)

	// create the file
	err = c.fs.WriteFile(filep, []byte(content), 0600)
	require.Nil(t, err)

	// read the file
	read, err := c.fs.ReadFile(filep)
	require.Nil(t, err)
	require.Equal(t, content, string(read))
}

func (c *conformance) TestWrite(t *testing.T, folder string, name string, data []byte, offset int64, write []byte, expected []byte) {
	require.NotNil(t, c.path)

	err := c.fs.MkdirAll(folder, 0666)
	require.Nil(t, err)

	full := c.path.Join(folder, name)

	f, err := c.fs.OpenFile(full, os.O_CREATE|os.O_RDWR, 0666)
	require.Nil(t, err)

	n, err := f.Write(data)
	require.Nil(t, err)
	require.Equal(t, len(data), n)

	n64, err := f.Seek(offset, io.SeekStart)
	require.Nil(t, err)
	require.Equal(t, offset, n64)

	n, err = f.Write(write)
	require.Nil(t, err)
	require.Equal(t, n, len(write))

	require.Nil(t, f.Close())

	content, err := c.fs.ReadFile(full)
	require.Nil(t, err)
	require.Equal(t, expected, content)
}

func (c *conformance) TestReadDir(t *testing.T, path string, files []file) {
	// both must equal
	err := c.fs.MkdirAll(path, 0666)
	require.Nil(t, err)

	fileNameMap := map[string]file{}
	for _, file := range files {
		filep := c.path.Join(path, file.name)
		content := file.content
		err = c.fs.WriteFile(filep, content, 0600)
		require.Nil(t, err)
		fileNameMap[file.name] = file
	}

	// list the files
	entries, err := c.fs.ReadDir(path)
	require.Nil(t, err)
	require.NotEmpty(t, entries)

	require.Equal(t, len(files), len(entries))

	// check the entry names and values
	for _, entry := range entries {
		require.Contains(t, fileNameMap, entry.Name())
	}
}

func (c *conformance) TestCanCreateFile(t *testing.T, path string, files []file) {
	err := c.fs.MkdirAll(path, 0666)
	require.Nil(t, err)

	// create the files
	for _, file := range files {
		filep := c.path.Join(path, file.name)
		f, err := c.fs.Create(filep)
		require.Nil(t, err)
		require.NotNil(t, file)
		io.Copy(f, bytes.NewBuffer(file.content))
		require.Nil(t, f.Close())
	}
}

func (c *conformance) TestCanWriteFile(t *testing.T, path string, files []file) {

	err := c.fs.MkdirAll(path, 0666)
	require.Nil(t, err)

	for _, file := range files {
		filep := c.path.Join(path, file.name)

		f, err := c.fs.Create(filep)
		require.Nil(t, err)
		require.NotNil(t, file)

		_, err = f.Write(file.content)
		require.Nil(t, err)
		require.Nil(t, f.Close())

		ofile, err := c.fs.Open(filep)
		require.Nil(t, err)

		stat, err := ofile.Stat()
		require.Nil(t, err)

		buf := make([]byte, stat.Size())
		_, err = ofile.Read(buf)
		require.Nil(t, err)
		require.Equal(t, file.content, buf)
		require.Nil(t, f.Close())
	}
}

func (c *conformance) TestWindowsWillNormalizePath(t *testing.T, folder string, file string) {

	err := c.fs.MkdirAll(folder, 0666)
	require.Nil(t, err)

	err = c.fs.WriteFile(c.path.Join(folder, file), []byte("content"), 0666)
	require.Nil(t, err)

	lower := strings.ToLower(c.path.Join(folder, file))
	ok, err := c.fs.Exists(lower)
	require.Nil(t, err)
	require.True(t, ok)
}

func (c *conformance) TestWindowsFileForwardAndBackwardSlash(t *testing.T, filePath string) {
	dir := c.path.Dir(filePath)
	err := c.fs.MkdirAll(dir, 0666)
	require.Nil(t, err)

	backPath := strings.ReplaceAll(filePath, "/", "\\")
	err = c.fs.WriteFile(backPath, []byte("test"), 0666)
	require.Nil(t, err)

	exists, err := c.fs.Exists(filePath)
	require.Nil(t, err)
	require.True(t, exists)
}

func (c *conformance) TestOpenFileFailsWhenNotExists(t *testing.T, folder string, file string) {
	err := c.fs.MkdirAll(folder, 0666)
	require.Nil(t, err)

	_, err = c.fs.OpenFile(file, os.O_RDONLY, 0666)
	require.NotNil(t, err)
}
