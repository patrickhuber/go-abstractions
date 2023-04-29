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
	windowsparse := []test{
		{path: "c:", fp: xfilepath.FilePath{
			Volume:   xfilepath.Volume{Drive: xfilepath.NullableString{Value: "c:", HasValue: true}},
			Absolute: false,
		}},
		{path: "c:/", fp: xfilepath.FilePath{
			Volume:   xfilepath.Volume{Drive: xfilepath.NullableString{Value: "c:", HasValue: true}},
			Segments: []string{""},
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
			Segments: []string{""},
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
			Segments: []string{""}, // trailing slash
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
		{path: "//abc//", fp: xfilepath.FilePath{
			Volume: xfilepath.Volume{
				Host:  xfilepath.NullableString{Value: "abc", HasValue: true},
				Share: xfilepath.NullableString{Value: "", HasValue: true},
			},
			Segments: []string{""},
			Absolute: true,
		}},
		{path: "//abc", fp: xfilepath.FilePath{
			Volume: xfilepath.Volume{
				Host: xfilepath.NullableString{Value: "abc", HasValue: true},
			},
			Absolute: true,
		}},
		{path: "///abc", fp: xfilepath.FilePath{
			Volume: xfilepath.Volume{
				Host:  xfilepath.NullableString{Value: "", HasValue: true},
				Share: xfilepath.NullableString{Value: "abc", HasValue: true},
			},
			Absolute: true,
		}},
		{path: "//abc//", fp: xfilepath.FilePath{
			Volume: xfilepath.Volume{
				Host:  xfilepath.NullableString{Value: "abc", HasValue: true},
				Share: xfilepath.NullableString{Value: "", HasValue: true},
			},
			Segments: []string{""},
			Absolute: true,
		}},
		{path: "///abc/", fp: xfilepath.FilePath{
			Volume: xfilepath.Volume{
				Host:  xfilepath.NullableString{Value: "", HasValue: true},
				Share: xfilepath.NullableString{Value: "abc", HasValue: true},
			},
			Absolute: true,
			Segments: []string{""},
		}},
		{path: "a/b", fp: xfilepath.FilePath{
			Segments: []string{"a", "b"},
			Absolute: false,
		}},
		{path: "a/b/", fp: xfilepath.FilePath{
			Segments: []string{"a", "b", ""},
			Absolute: false,
		}},
		{path: "a/", fp: xfilepath.FilePath{
			Segments: []string{"a", ""},
			Absolute: false,
		}},
		{path: "a", fp: xfilepath.FilePath{
			Segments: []string{"a"},
			Absolute: false,
		}},
		{path: "", fp: xfilepath.FilePath{
			Absolute: false,
		}},
		// {path: `\a`, fp: xfilepath.FilePath{
		// 	Absolute: false,
		// 	Segments: []string{"a"},
		// }},
	}
	linuxparse := []test{
		{path: "/", fp: xfilepath.FilePath{
			Absolute: true,
		}},
	}
	run := func(tests []test, name string, plat platform.Platform) {
		for _, test := range tests {
			parser := xfilepath.NewParserWithPlatform(plat)
			actual, err := parser.Parse(test.path)
			require.Nil(t, err)
			require.Equal(t, test.fp, actual, "unable to parse path '%s'", test.path)
		}
	}
	run(windowsparse, "windowsparse", platform.Windows)
	run(linuxparse, "linuxparse", platform.Linux)

}
