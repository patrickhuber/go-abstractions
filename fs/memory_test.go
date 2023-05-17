package fs_test

import (
	"regexp"
	"testing"

	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/stretchr/testify/require"
)

func TestMkdirCreatesRootUnix(t *testing.T) {
	path := "/"
	parser := filepath.NewParserWithPlatform(platform.Linux)
	processor := filepath.NewProcessor(filepath.WithParser(parser))
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
	processor := filepath.NewProcessorWithPlatform(platform.Linux)
	f := fs.NewMemory(fs.WithProcessor(processor))
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
