package fs_test

import (
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/stretchr/testify/require"
)

func TestMkdirCreatesRootUnix(t *testing.T) {
	path := "/"
	processor := filepath.NewProcessorWithPlatform(platform.Linux)
	f := fs.NewMemory(fs.WithProcessor(processor))
	err := f.Mkdir(path, 0666)
	require.Nil(t, err)
	ok, err := f.Exists(path)
	require.Nil(t, err)
	require.True(t, ok)
}

func TestMkdirFailsWhenRootNotExists(t *testing.T) {
	path := "/test"
	processor := filepath.NewProcessorWithPlatform(platform.Linux)
	f := fs.NewMemory(fs.WithProcessor(processor))
	err := f.Mkdir(path, 0666)
	require.NotNil(t, err)
}

func TestMkdirAllCreatesAllDirectories(t *testing.T) {
	path := "/gran/parent/child"
	f, _ := setup(platform.Linux)
	err := f.MkdirAll(path, 0666)
	require.Nil(t, err)
	paths := []string{
		"/",
		"/gran",
		"/gran/parent",
		"/gran/parent/child",
	}
	for _, p := range paths {
		ok, err := f.Exists(p)
		require.Nil(t, err)
		require.True(t, ok, "%s does not exist", p)
	}
}

func TestWriteFile(t *testing.T) {
	path := "/gran/parent/child"
	processor := filepath.NewProcessorWithPlatform(platform.Linux)
	filep := processor.Join(path, "file.txt")
	f := fs.NewMemory(fs.WithProcessor(processor))
	err := f.MkdirAll(path, 0666)
	require.Nil(t, err)

	// create the file
	err = f.WriteFile(filep, []byte("file"), 0600)
	require.Nil(t, err)

	// read the file
	content, err := f.ReadFile(filep)
	require.Nil(t, err)
	require.Equal(t, "file", string(content))
}

func TestWrite(t *testing.T) {
	folder := "/gran/parent/child"
	fs, path := setup(platform.Linux)
	type test struct {
		name     string
		data     []byte
		offset   int64
		write    []byte
		expected string
	}

	tests := []test{
		{"grow.txt", []byte("this is test data"), 4, []byte(" more data than expected"), "this is more data than expected"},
		{"less.txt", []byte("this is test data"), 8, []byte("also"), "this is also data"},
		{"end.txt", []byte("this is test data"), 13, []byte("info"), "this is test info"},
	}

	err := fs.MkdirAll(folder, 0666)
	require.Nil(t, err)

	for _, test := range tests {
		full := path.Join(folder, test.name)

		f, err := fs.OpenFile(full, os.O_CREATE|os.O_RDWR, 0666)
		require.Nil(t, err)

		n, err := f.Write(test.data)
		require.Nil(t, err)
		require.Equal(t, len(test.data), n)

		n64, err := f.Seek(test.offset, io.SeekStart)
		require.Nil(t, err)
		require.Equal(t, test.offset, n64)

		n, err = f.Write(test.write)
		require.Nil(t, err)
		require.Equal(t, n, len(test.write))

		require.Nil(t, f.Close())
	}
}

func TestReadDir(t *testing.T) {
	path := "/gran/parent/child"
	processor := filepath.NewProcessorWithPlatform(platform.Linux)
	files := []string{"one.txt", "two.txt", "three.txt"}
	f := fs.NewMemory(fs.WithProcessor(processor))
	err := f.MkdirAll(path, 0666)
	require.Nil(t, err)

	// write the files
	for _, file := range files {
		filep := processor.Join(path, file)
		err = f.WriteFile(filep, []byte(file), 0600)
		require.Nil(t, err)
	}

	// list the files
	entries, err := f.ReadDir(path)
	require.Nil(t, err)
	require.NotEmpty(t, entries)
	require.Equal(t, len(files), len(entries))

	// check the entry names and values
	nameMatch, err := regexp.Compile(`^\w+.txt$`)
	require.Nil(t, err)
	for _, entry := range entries {
		matched := nameMatch.MatchString(entry.Name())
		require.True(t, matched, "name %s is not correct", entry.Name())
	}
}

func TestCanCreateFile(t *testing.T) {
	path := "/gran/parent/child"
	processor := filepath.NewProcessorWithPlatform(platform.Linux)
	files := []string{"one.txt", "two.txt", "three.txt"}
	f := fs.NewMemory(fs.WithProcessor(processor))
	err := f.MkdirAll(path, 0666)
	require.Nil(t, err)

	// write the files
	for _, file := range files {
		filep := processor.Join(path, file)
		file, err := f.Create(filep)
		require.Nil(t, err)
		require.NotNil(t, file)
		require.Nil(t, file.Close())
	}
}

func TestCanWriteFile(t *testing.T) {
	f, processor := setup(platform.Linux)
	path := "/gran/parent/child"
	files := []string{"test.txt"}
	err := f.MkdirAll(path, 0666)
	require.Nil(t, err)

	for _, file := range files {
		filep := processor.Join(path, file)
		file, err := f.Create(filep)
		require.Nil(t, err)
		require.NotNil(t, file)
		_, err = file.Write([]byte("test"))
		require.Nil(t, err)
		require.Nil(t, file.Close())
		ofile, err := f.Open(filep)
		require.Nil(t, err)
		buf := []byte("    ")
		_, err = ofile.Read(buf)
		require.Nil(t, err)
		require.Equal(t, "test", string(buf))
		require.Nil(t, file.Close())
	}
}

func TestWindowsWillNormalizePath(t *testing.T) {
	fs, path := setup(platform.Windows)
	folder := `c:\ProgramData\fake\folder`
	file := `test.txt`

	err := fs.MkdirAll(folder, 0666)
	require.Nil(t, err)

	err = fs.WriteFile(path.Join(folder, file), []byte("content"), 0666)
	require.Nil(t, err)

	// lower case should not matter
	lower := strings.ToLower(path.Join(folder, file))
	ok, err := fs.Exists(lower)
	require.Nil(t, err)
	require.True(t, ok)

	// forward slashes should be the same as backward slashes
	forward := strings.ReplaceAll(path.Join(folder, file), `\`, "/")
	ok, err = fs.Exists(forward)
	require.Nil(t, err)
	require.True(t, ok)
}

func setup(plat platform.Platform) (fs.FS, filepath.Processor) {
	processor := filepath.NewProcessorWithPlatform(plat)
	fs := fs.NewMemory(fs.WithProcessor(processor))
	return fs, processor
}
