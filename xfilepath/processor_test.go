package xfilepath_test

import (
	"strings"
	"testing"

	"github.com/patrickhuber/go-xplat/xfilepath"
	"github.com/stretchr/testify/require"
)

func TestJoin(t *testing.T) {
	type test struct {
		elements []string
		sep      xfilepath.PathSeparator
		result   string
	}

	tests := []test{
		{
			[]string{"a", "b", "c"},
			xfilepath.ForwardSlash,
			"a/b/c",
		},
		{
			[]string{"a", "b/c"},
			xfilepath.ForwardSlash,
			"a/b/c",
		},
		{
			[]string{"a/b", "c"},
			xfilepath.ForwardSlash,
			"a/b/c",
		},
		{
			[]string{"a/b", "/c"},
			xfilepath.ForwardSlash,
			"a/b/c",
		},
		{
			[]string{"/a/b", "/c"},
			xfilepath.ForwardSlash,
			"/a/b/c",
		},
		{
			[]string{`c:\`, `a\b`, `c`},
			xfilepath.BackwardSlash,
			`c:\a\b\c`,
		},
	}

	for _, test := range tests {
		processor := xfilepath.NewProcessorWith(test.sep)
		actual := processor.Join(test.elements...)
		require.Equal(t, test.result, actual)
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
			"//host/share/gran/parent/child",
			xfilepath.ForwardSlash,
			"//host/share",
		},
		{
			// UNC backward slash
			`\\host\share\gran\parent\child`,
			xfilepath.BackwardSlash,
			`\\host\share`,
		},
		{
			// Unix Path
			"/gran/parent/child",
			xfilepath.ForwardSlash,
			"/",
		},
		{
			// Windows Path
			`c:\gran\parent\child`,
			xfilepath.BackwardSlash,
			`c:`,
		},
	}

	for _, test := range tests {
		processor := xfilepath.NewProcessorWith(test.sep)
		actual := processor.Root(test.path)
		require.Equal(t, test.expected, actual)
	}
}

func TestRel(t *testing.T) {
	type test struct {
		source   string
		target   string
		expected string
	}
	tests := []test{
		{"a/b", "a/b", "."},
		{"a/b/.", "a/b", "."},
		{"a/b", "a/b/.", "."},
		{"./a/b", "a/b", "."},
		{"a/b", "./a/b", "."},
		{"ab/cd", "ab/cde", "../cde"},
		{"ab/cd", "ab/c", "../c"},
		{"a/b", "a/b/c/d", "c/d"},
		{"a/b", "a/b/../c", "../c"},
		{"a/b/../c", "a/b", "../b"},
		{"a/b/c", "a/c/d", "../../c/d"},
		{"a/b", "c/d", "../../c/d"},
		{"a/b/c/d", "a/b", "../.."},
		{"a/b/c/d", "a/b/", "../.."},
		{"a/b/c/d/", "a/b", "../.."},
		{"a/b/c/d/", "a/b/", "../.."},
		{"../../a/b", "../../a/b/c/d", "c/d"},
		{"/a/b", "/a/b", "."},
		{"/a/b/.", "/a/b", "."},
		{"/a/b", "/a/b/.", "."},
		{"/ab/cd", "/ab/cde", "../cde"},
		{"/ab/cd", "/ab/c", "../c"},
		{"/a/b", "/a/b/c/d", "c/d"},
		{"/a/b", "/a/b/../c", "../c"},
		{"/a/b/../c", "/a/b", "../b"},
		{"/a/b/c", "/a/c/d", "../../c/d"},
		{"/a/b", "/c/d", "../../c/d"},
		{"/a/b/c/d", "/a/b", "../.."},
		{"/a/b/c/d", "/a/b/", "../.."},
		{"/a/b/c/d/", "/a/b", "../.."},
		{"/a/b/c/d/", "/a/b/", "../.."},
		{"/../../a/b", "/../../a/b/c/d", "c/d"},
		{".", "a/b", "a/b"},
		{".", "..", ".."},
	}
	for i, test := range tests {
		processor := xfilepath.NewProcessorWith(xfilepath.ForwardSlash)
		actual, err := processor.Rel(test.source, test.target)
		require.Nil(t, err)
		require.Equal(t, test.expected, actual, "failed on test %d", i)
	}
}

func TestClean(t *testing.T) {
	type test struct {
		path     string
		expected string
	}
	tests := []test{
		// Already clean
		{"abc", "abc"},
		{"abc/def", "abc/def"},
		{"a/b/c", "a/b/c"},
		{".", "."},
		{"..", ".."},
		{"../..", "../.."},
		{"../../abc", "../../abc"},
		{"/abc", "/abc"},
		{"/", "/"},

		// Empty is current dir
		{"", "."},

		// Remove trailing slash
		{"abc/", "abc"},
		{"abc/def/", "abc/def"},
		{"a/b/c/", "a/b/c"},
		{"./", "."},
		{"../", ".."},
		{"../../", "../.."},
		{"/abc/", "/abc"},

		// Remove doubled slash
		{"abc//def//ghi", "abc/def/ghi"},
		{"abc//", "abc"},

		// Remove . elements
		{"abc/./def", "abc/def"},
		{"/./abc/def", "/abc/def"},
		{"abc/.", "abc"},

		// Remove .. elements
		{"abc/def/ghi/../jkl", "abc/def/jkl"},
		{"abc/def/../ghi/../jkl", "abc/jkl"},
		{"abc/def/..", "abc"},
		{"abc/def/../..", "."},
		{"/abc/def/../..", "/"},
		{"abc/def/../../..", ".."},
		{"/abc/def/../../..", "/"},
		{"abc/def/../../../ghi/jkl/../../../mno", "../../mno"},
		{"/../abc", "/abc"},

		// Combinations
		{"abc/./../def", "def"},
		{"abc//./../def", "def"},
		{"abc/../../././../def", "../../def"},

		// Remove leading doubled slash (these will be parsed as UNC in windows)
		// {"//abc", "/abc"},
		// {"///abc", "/abc"},
		// {"//abc//", "/abc"},

		// windows clean
		{`c:`, `c:.`},
		{`c:\`, `c:\`},
		{`c:\abc`, `c:\abc`},
		{`c:abc\..\..\.\.\..\def`, `c:..\..\def`},
		{`c:\abc\def\..\..`, `c:\`},
		{`c:\..\abc`, `c:\abc`},
		{`c:..\abc`, `c:..\abc`},
		{`\`, `\`},
		{`/`, `\`},
		{`\\i\..\c$`, `\\i\..\c$`},
		{`\\i\..\i\c$`, `\\i\..\i\c$`},
		{`\\i\..\I\c$`, `\\i\..\I\c$`},
		{`\\host\share\foo\..\bar`, `\\host\share\bar`},
		{`//host/share/foo/../baz`, `\\host\share\baz`},
		{`\\host\share\foo\..\..\..\..\bar`, `\\host\share\bar`},
		{`\\.\C:\a\..\..\..\..\bar`, `\\.\C:\bar`},
		{`\\.\C:\\\\a`, `\\.\C:\a`},
		{`\\a\b\..\c`, `\\a\b\c`},
		{`\\a\b`, `\\a\b`},
		{`.\c:`, `.\c:`},
		{`.\c:\foo`, `.\c:\foo`},
		{`.\c:foo`, `.\c:foo`},
		{`//abc`, `\\abc`},
		{`///abc`, `\\\abc`},
		{`//abc//`, `\\abc\\`},

		// Don't allow cleaning to move an element with a colon to the start of the path.
		{`a/../c:`, `.\c:`},
		{`a\..\c:`, `.\c:`},
		{`a/../c:/a`, `.\c:\a`},
		{`a/../../c:`, `..\c:`},
		{`foo:bar`, `foo:bar`},
	}

	for i, test := range tests {
		sep := xfilepath.BackwardSlash
		if strings.Contains(test.expected, "/") {
			sep = xfilepath.ForwardSlash
		}
		processor := xfilepath.NewProcessorWith(sep)
		actual := processor.Clean(test.path)
		require.Equal(t, test.expected, actual, "failed on test %d", i)
	}
}
