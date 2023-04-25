package xfilepath_test

import (
	"testing"

	"github.com/patrickhuber/go-xplat/xfilepath"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	type test struct {
		fp       xfilepath.FilePath
		sep      xfilepath.PathSeparator
		expected string
	}

	tests := []test{
		{
			// UNC share forward slash
			fp: xfilepath.FilePath{
				Volume: xfilepath.Volume{
					Host:  xfilepath.NullableString{Value: "host", HasValue: true},
					Share: xfilepath.NullableString{Value: "share", HasValue: true},
				},
				Segments: []string{"gran", "parent", "child"},
				Absolute: true,
			},
			sep:      xfilepath.ForwardSlash,
			expected: "//host/share/gran/parent/child",
		},
		{
			// UNC share backward slash
			fp: xfilepath.FilePath{
				Volume: xfilepath.Volume{
					Host:  xfilepath.NullableString{Value: "host", HasValue: true},
					Share: xfilepath.NullableString{Value: "share", HasValue: true},
				},
				Segments: []string{"gran", "parent", "child"},
				Absolute: true,
			},
			sep:      xfilepath.BackwardSlash,
			expected: `\\host\share\gran\parent\child`,
		},
		{
			// UNC share only root
			fp: xfilepath.FilePath{
				Volume: xfilepath.Volume{
					Host:  xfilepath.NullableString{Value: "host", HasValue: true},
					Share: xfilepath.NullableString{Value: "share", HasValue: true},
				},
				Segments: []string{},
				Absolute: true,
			},
			sep:      xfilepath.BackwardSlash,
			expected: `\\host\share`,
		},
		{
			// Unix Path
			fp: xfilepath.FilePath{
				Segments: []string{"gran", "parent", "child"},
				Absolute: true,
			},
			sep:      xfilepath.ForwardSlash,
			expected: "/gran/parent/child",
		},
		{
			// Windows Path
			fp: xfilepath.FilePath{
				Volume: xfilepath.Volume{
					Drive: xfilepath.NullableString{Value: `c:`, HasValue: true},
				},
				Segments: []string{"gran", "parent", "child"},
				Absolute: true,
			},
			sep:      xfilepath.BackwardSlash,
			expected: `c:\gran\parent\child`,
		},
		{
			// relative
			fp: xfilepath.FilePath{
				Segments: []string{"gran", "parent", "child"},
				Absolute: false,
			},
			sep:      xfilepath.ForwardSlash,
			expected: "gran/parent/child",
		},
		{
			// root unix path
			fp: xfilepath.FilePath{
				Absolute: true,
			},
			sep:      xfilepath.ForwardSlash,
			expected: "/",
		},
		{
			// root windows path
			fp: xfilepath.FilePath{
				Volume: xfilepath.Volume{
					Drive: xfilepath.NullableString{Value: `c:`, HasValue: true},
				},
				Absolute: true,
			},
			sep:      xfilepath.BackwardSlash,
			expected: `c:\`,
		},
	}

	for _, test := range tests {
		processor := xfilepath.NewProcessor(xfilepath.WithSeparator(test.sep))
		actual := test.fp.String(processor.Separator())
		require.Equal(t, test.expected, actual)
	}
}

func TestVolumeName(t *testing.T) {
	type test struct {
		path     string
		expected string
		sep      xfilepath.PathSeparator
	}

	tests := []test{
		{
			"//host/share/gran/parent/child",
			"//host/share",
			xfilepath.ForwardSlash,
		},
		{
			`\\host\share\gran\parent\child`,
			`\\host\share`,
			xfilepath.BackwardSlash,
		},
		{
			"/gran/parent/child",
			"",
			xfilepath.ForwardSlash,
		},
		{
			// Windows Path
			`c:\gran\parent\child`,
			`c:`,
			xfilepath.BackwardSlash,
		},
	}

	for _, test := range tests {
		processor := xfilepath.NewProcessor(xfilepath.WithSeparator(test.sep))
		actual := processor.VolumeName(test.path)
		require.Equal(t, test.expected, actual)
	}
}
