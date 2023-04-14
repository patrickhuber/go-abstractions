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
					Host:  "host",
					Share: "share",
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
					Host:  "host",
					Share: "share",
				},
				Segments: []string{"gran", "parent", "child"},
				Absolute: true,
			},
			sep:      xfilepath.BackwardSlash,
			expected: `\\host\share\gran\parent\child`,
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
					Drive: `c:`,
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
			fp: xfilepath.FilePath{
				Absolute: true,
			},
			sep:      xfilepath.ForwardSlash,
			expected: "/",
		},
	}

	for _, test := range tests {
		actual := test.fp.String(test.sep)
		require.Equal(t, test.expected, actual)
	}
}

func TestRoot(t *testing.T) {
	type test struct {
		path     string
		sep      xfilepath.PathSeparator
		expected string
	}

	tests := []test{
		{
			// UNC forward slash
			path:     "//host/share/gran/parent/child",
			sep:      xfilepath.ForwardSlash,
			expected: "//host/share",
		},
		{
			// UNC backward slash
			path:     `\\host\share\gran\parent\child`,
			sep:      xfilepath.BackwardSlash,
			expected: `\\host\share`,
		},
		{
			// Unix Path
			path:     "/gran/parent/child",
			sep:      xfilepath.ForwardSlash,
			expected: "/",
		},
		{
			// Windows Path
			path:     `c:\gran\parent\child`,
			sep:      xfilepath.BackwardSlash,
			expected: `c:`,
		},
	}

	for _, test := range tests {
		actual := xfilepath.Root(test.sep, test.path)
		require.Equal(t, test.expected, actual)
	}
}

func TestVolumeName(t *testing.T) {
	type test struct {
		fp       xfilepath.FilePath
		sep      xfilepath.PathSeparator
		expected string
	}

	tests := []test{
		{
			// UNC forward slash
			fp: xfilepath.FilePath{
				Volume: xfilepath.Volume{
					Host:  "host",
					Share: "share",
				},
				Segments: []string{"gran", "parent", "child"},
				Absolute: true,
			},
			sep:      xfilepath.ForwardSlash,
			expected: "//host/share",
		},
		{
			// UNC backward slash
			fp: xfilepath.FilePath{
				Volume: xfilepath.Volume{
					Host:  "host",
					Share: "share",
				},
				Segments: []string{"gran", "parent", "child"},
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
			expected: "",
		},
		{
			// Windows Path
			fp: xfilepath.FilePath{
				Volume: xfilepath.Volume{
					Drive: `c:`,
				},
				Segments: []string{"gran", "parent", "child"},
				Absolute: true,
			},
			sep:      xfilepath.BackwardSlash,
			expected: `c:`,
		},
	}

	for _, test := range tests {
		actual := test.fp.VolumeName(test.sep)
		require.Equal(t, test.expected, actual)
	}
}

func TestJoin(t *testing.T) {
	type test struct {
		elements []string
		sep      xfilepath.PathSeparator
		result   string
	}

	tests := []test{
		{
			elements: []string{"a", "b", "c"},
			sep:      xfilepath.ForwardSlash,
			result:   "a/b/c",
		},
		{
			elements: []string{"a", "b/c"},
			sep:      xfilepath.ForwardSlash,
			result:   "a/b/c",
		},
		{
			elements: []string{"a/b", "c"},
			sep:      xfilepath.ForwardSlash,
			result:   "a/b/c",
		},
		{
			elements: []string{"a/b", "/c"},
			sep:      xfilepath.ForwardSlash,
			result:   "a/b/c",
		},
		{
			elements: []string{"/a/b", "/c"},
			sep:      xfilepath.ForwardSlash,
			result:   "/a/b/c",
		},
		{
			elements: []string{`c:\`, `a\b`, `c`},
			sep:      xfilepath.BackwardSlash,
			result:   `c:\a\b\c`,
		},
	}

	for _, test := range tests {
		actual := xfilepath.Join(test.sep, test.elements...)
		require.Equal(t, test.result, actual)
	}
}
