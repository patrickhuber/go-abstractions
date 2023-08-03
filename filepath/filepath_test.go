package filepath_test

import (
	"testing"

	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	type test struct {
		fp       filepath.FilePath
		sep      filepath.PathSeparator
		expected string
	}

	tests := []test{
		{
			// UNC share forward slash
			fp: filepath.FilePath{
				Volume: filepath.Volume{
					Host:  filepath.NullableString{Value: "host", HasValue: true},
					Share: filepath.NullableString{Value: "share", HasValue: true},
				},
				Segments: []string{"gran", "parent", "child"},
				Absolute: true,
			},
			sep:      filepath.ForwardSlash,
			expected: "//host/share/gran/parent/child",
		},
		{
			// UNC share backward slash
			fp: filepath.FilePath{
				Volume: filepath.Volume{
					Host:  filepath.NullableString{Value: "host", HasValue: true},
					Share: filepath.NullableString{Value: "share", HasValue: true},
				},
				Segments: []string{"gran", "parent", "child"},
				Absolute: true,
			},
			sep:      filepath.BackwardSlash,
			expected: `\\host\share\gran\parent\child`,
		},
		{
			// UNC share only root
			fp: filepath.FilePath{
				Volume: filepath.Volume{
					Host:  filepath.NullableString{Value: "host", HasValue: true},
					Share: filepath.NullableString{Value: "share", HasValue: true},
				},
				Segments: []string{},
				Absolute: true,
			},
			sep:      filepath.BackwardSlash,
			expected: `\\host\share`,
		},
		{
			// UNC share without share name and trailing slash
			fp: filepath.FilePath{
				Volume: filepath.Volume{
					Host:  filepath.NullableString{Value: "abc", HasValue: true},
					Share: filepath.NullableString{Value: "", HasValue: true},
				},
				Segments: []string{""},
				Absolute: true,
			},
			sep:      filepath.BackwardSlash,
			expected: `\\abc\\`,
		},
		{
			// UNC share without share name and trailing slash
			fp: filepath.FilePath{
				Volume: filepath.Volume{
					Host: filepath.NullableString{Value: "abc", HasValue: true},
				},
				Absolute: true,
			},
			sep:      filepath.BackwardSlash,
			expected: `\\abc`,
		},
		{
			// Unix Path
			fp: filepath.FilePath{
				Segments: []string{"gran", "parent", "child"},
				Absolute: true,
			},
			sep:      filepath.ForwardSlash,
			expected: "/gran/parent/child",
		},
		{
			// Windows Path
			fp: filepath.FilePath{
				Volume: filepath.Volume{
					Drive: filepath.NullableString{Value: `c:`, HasValue: true},
				},
				Segments: []string{"gran", "parent", "child"},
				Absolute: true,
			},
			sep:      filepath.BackwardSlash,
			expected: `c:\gran\parent\child`,
		},
		{
			// relative
			fp: filepath.FilePath{
				Segments: []string{"gran", "parent", "child"},
				Absolute: false,
			},
			sep:      filepath.ForwardSlash,
			expected: "gran/parent/child",
		},
		{
			// root unix path
			fp: filepath.FilePath{
				Absolute: true,
			},
			sep:      filepath.ForwardSlash,
			expected: "/",
		},
		{
			// root windows path
			fp: filepath.FilePath{
				Volume: filepath.Volume{
					Drive: filepath.NullableString{Value: `c:`, HasValue: true},
				},
				Absolute: true,
			},
			sep:      filepath.BackwardSlash,
			expected: `c:\`,
		},
	}

	for i, test := range tests {
		processor := filepath.NewProcessor()
		processor.Separator = test.sep
		actual := test.fp.String(processor.Separator)
		require.Equal(t, test.expected, actual, "failed test at [%d]", i)
	}
}

func TestVolumeName(t *testing.T) {
	type test struct {
		path     string
		expected string
		sep      filepath.PathSeparator
	}

	tests := []test{
		{
			"//host/share/gran/parent/child",
			"//host/share",
			filepath.ForwardSlash,
		},
		{
			`\\host\share\gran\parent\child`,
			`\\host\share`,
			filepath.BackwardSlash,
		},
		{
			"/gran/parent/child",
			"",
			filepath.ForwardSlash,
		},
		{
			// Windows Path
			`c:\gran\parent\child`,
			`c:`,
			filepath.BackwardSlash,
		},
	}

	for _, test := range tests {
		processor := filepath.NewProcessor()
		processor.Separator = test.sep
		actual := processor.VolumeName(test.path)
		require.Equal(t, test.expected, actual)
	}
}
