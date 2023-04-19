package xfilepath_test

import (
	"testing"

	"github.com/patrickhuber/go-xplat/xfilepath"
	"github.com/stretchr/testify/require"
)

func TestCanParse(t *testing.T) {
	type test struct {
		path string
		fp   xfilepath.FilePath
	}

	tests := []test{
		{path: "c:", fp: xfilepath.FilePath{
			Volume:   xfilepath.Volume{Drive: "c:"},
			Absolute: true,
		}},
		{path: "c:/", fp: xfilepath.FilePath{
			Volume:   xfilepath.Volume{Drive: "c:"},
			Absolute: true,
			Trailing: true,
		}},
		{path: "c:/foo", fp: xfilepath.FilePath{
			Volume:   xfilepath.Volume{Drive: "c:"},
			Segments: []string{"foo"},
			Absolute: true,
		}},
		{path: "c:/foo/bar", fp: xfilepath.FilePath{
			Volume:   xfilepath.Volume{Drive: "c:"},
			Segments: []string{"foo", "bar"},
			Absolute: true,
		}},
		{path: "//host/share", fp: xfilepath.FilePath{
			Volume:   xfilepath.Volume{Host: "host", Share: "share"},
			Absolute: true,
		}},
		{path: "//host/share/", fp: xfilepath.FilePath{
			Volume:   xfilepath.Volume{Host: "host", Share: "share"},
			Absolute: true,
			Trailing: true,
		}},
		{path: "//host/share/foo", fp: xfilepath.FilePath{
			Volume:   xfilepath.Volume{Host: "host", Share: "share"},
			Segments: []string{"foo"},
			Absolute: true,
		}},
		{path: `\\host\share`, fp: xfilepath.FilePath{
			Volume:   xfilepath.Volume{Host: "host", Share: "share"},
			Absolute: true,
		}},
		{path: `\\host\share\`, fp: xfilepath.FilePath{
			Volume:   xfilepath.Volume{Host: "host", Share: "share"},
			Absolute: true,
			Trailing: true,
		}},
		{path: `\\host\share\foo`, fp: xfilepath.FilePath{
			Volume:   xfilepath.Volume{Host: "host", Share: "share"},
			Segments: []string{"foo"},
			Absolute: true,
		}},
		{path: `//./NUL`, fp: xfilepath.FilePath{
			Volume:   xfilepath.Volume{Host: ".", Share: "NUL"},
			Absolute: true,
		}},
		{path: `//?/NUL`, fp: xfilepath.FilePath{
			Volume:   xfilepath.Volume{Host: "?", Share: "NUL"},
			Absolute: true,
		}},

		{path: "/", fp: xfilepath.FilePath{
			Absolute: true,
			Trailing: true,
		}},
		{path: "a/b", fp: xfilepath.FilePath{
			Segments: []string{"a", "b"},
			Absolute: false,
		}},
		{path: "a/b/", fp: xfilepath.FilePath{
			Segments: []string{"a", "b"},
			Absolute: false,
			Trailing: true,
		}},
		{path: "a/", fp: xfilepath.FilePath{
			Segments: []string{"a"},
			Absolute: false,
			Trailing: true,
		}},
		{path: "a", fp: xfilepath.FilePath{
			Segments: []string{"a"},
			Absolute: false,
		}},
	}

	for _, test := range tests {
		result, err := xfilepath.Parse(test.path)
		require.Nil(t, err)
		require.Equal(t, test.fp, result, "unable to parse path '%s'", test.path)
	}
}
