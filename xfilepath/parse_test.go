package xfilepath_test

import (
	"testing"

	"github.com/patrickhuber/go-xplat/platform"
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
			Volume:   xfilepath.Volume{Drive: xfilepath.NullableString{Value: "c:", HasValue: true}},
			Absolute: false,
		}},
		{path: "c:/", fp: xfilepath.FilePath{
			Volume:   xfilepath.Volume{Drive: xfilepath.NullableString{Value: "c:", HasValue: true}},
			Absolute: true,
		}},
		{path: "c:/foo", fp: xfilepath.FilePath{
			Volume:   xfilepath.Volume{Drive: xfilepath.NullableString{Value: "c:", HasValue: true}},
			Segments: []string{"foo"},
			Absolute: true,
		}},
		{path: "c:/foo/bar", fp: xfilepath.FilePath{
			Volume:   xfilepath.Volume{Drive: xfilepath.NullableString{Value: "c:", HasValue: true}},
			Segments: []string{"foo", "bar"},
			Absolute: true,
		}},
		{path: "//host/share", fp: xfilepath.FilePath{
			Volume: xfilepath.Volume{
				Host:  xfilepath.NullableString{Value: "host", HasValue: true},
				Share: xfilepath.NullableString{Value: "share", HasValue: true},
			},
			Absolute: true,
		}},
		{path: "//host/share/", fp: xfilepath.FilePath{
			Volume: xfilepath.Volume{
				Host:  xfilepath.NullableString{Value: "host", HasValue: true},
				Share: xfilepath.NullableString{Value: "share", HasValue: true},
			},
			Absolute: true,
		}},
		{path: "//host/share/foo", fp: xfilepath.FilePath{
			Volume: xfilepath.Volume{
				Host:  xfilepath.NullableString{Value: "host", HasValue: true},
				Share: xfilepath.NullableString{Value: "share", HasValue: true},
			},
			Segments: []string{"foo"},
			Absolute: true,
		}},
		{path: `\\host\share`, fp: xfilepath.FilePath{
			Volume: xfilepath.Volume{
				Host:  xfilepath.NullableString{Value: "host", HasValue: true},
				Share: xfilepath.NullableString{Value: "share", HasValue: true},
			},
			Absolute: true,
		}},
		{path: `\\host\share\`, fp: xfilepath.FilePath{
			Volume: xfilepath.Volume{
				Host:  xfilepath.NullableString{Value: "host", HasValue: true},
				Share: xfilepath.NullableString{Value: "share", HasValue: true},
			},
			Absolute: true,
		}},
		{path: `\\host\share\foo`, fp: xfilepath.FilePath{
			Volume: xfilepath.Volume{
				Host:  xfilepath.NullableString{Value: "host", HasValue: true},
				Share: xfilepath.NullableString{Value: "share", HasValue: true},
			},
			Segments: []string{"foo"},
			Absolute: true,
		}},
		{path: `//./NUL`, fp: xfilepath.FilePath{
			Volume: xfilepath.Volume{
				Host:  xfilepath.NullableString{Value: ".", HasValue: true},
				Share: xfilepath.NullableString{Value: "NUL", HasValue: true},
			},
			Absolute: true,
		}},
		{path: `//?/NUL`, fp: xfilepath.FilePath{
			Volume: xfilepath.Volume{
				Host:  xfilepath.NullableString{Value: "?", HasValue: true},
				Share: xfilepath.NullableString{Value: "NUL", HasValue: true},
			},
			Absolute: true,
		}},

		{path: "/", fp: xfilepath.FilePath{
			Absolute: true,
		}},
		{path: "a/b", fp: xfilepath.FilePath{
			Segments: []string{"a", "b"},
			Absolute: false,
		}},
		{path: "a/b/", fp: xfilepath.FilePath{
			Segments: []string{"a", "b"},
			Absolute: false,
		}},
		{path: "a/", fp: xfilepath.FilePath{
			Segments: []string{"a"},
			Absolute: false,
		}},
		{path: "a", fp: xfilepath.FilePath{
			Segments: []string{"a"},
			Absolute: false,
		}},
		{path: "", fp: xfilepath.FilePath{
			Absolute: false,
		}},
	}

	for _, test := range tests {
		parser := xfilepath.NewParserWithPlatform(platform.Windows)
		result, err := parser.Parse(test.path)
		require.Nil(t, err)
		require.Equal(t, test.fp, result, "unable to parse path '%s'", test.path)
	}
}
