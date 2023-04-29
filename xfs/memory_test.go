package xfs_test

import (
	"testing"

	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/go-xplat/xfilepath"
	"github.com/patrickhuber/go-xplat/xfs"
	"github.com/stretchr/testify/require"
)

func TestMkdirCreatesRootUnix(t *testing.T) {
	path := "/"
	parser := xfilepath.NewParserWithPlatform(platform.Linux)
	processor := xfilepath.NewProcessor(xfilepath.WithParser(parser))
	f := xfs.NewMemory(xfs.WithProcessor(processor))
	err := f.Mkdir(path, 0666)
	require.Nil(t, err)
	ok, err := f.Exists(path)
	require.Nil(t, err)
	require.True(t, ok)
}

func TestMkdirFailsWhenRootNotExists(t *testing.T) {
	path := "/test"
	processor := xfilepath.NewProcessorWithPlatform(platform.Linux)
	f := xfs.NewMemory(xfs.WithProcessor(processor))
	err := f.Mkdir(path, 0666)
	require.NotNil(t, err)
}

func TestMkdirAllCreatesAllDirectories(t *testing.T) {
	path := "/gran/parent/child"
	processor := xfilepath.NewProcessorWithPlatform(platform.Linux)
	f := xfs.NewMemory(xfs.WithProcessor(processor))
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
