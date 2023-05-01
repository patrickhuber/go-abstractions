package fs_test

import (
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
