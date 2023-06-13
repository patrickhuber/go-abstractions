package fs_test

import (
	"testing"

	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/platform"
)

func TestMemoryMkdirCreatesRootUnix(t *testing.T) {
	NewConformanceWithPath(setupMemory(platform.Linux)).TestMkdirCreatesRoot(t, "/")
}

func TestMemoryMkdirFailsWhenRootNotExists(t *testing.T) {
	NewConformanceWithPath(setupMemory(platform.Linux)).TestMkdirFailsWhenRootNotExists(t, "/test")
}

func TestMemoryMkdirAllCreatesAllDirectories(t *testing.T) {
	NewConformanceWithPath(setupMemory(platform.Linux)).TestMkdirAllCreatesAllDirectories(t, "/gran/parent/child", []string{
		"/",
		"/gran",
		"/gran/parent",
		"/gran/parent/child",
	})
}

func TestMemoryWriteFile(t *testing.T) {
	NewConformanceWithPath(setupMemory(platform.Linux)).TestWriteFile(t, "/gran/parent/child", "file.txt", "file")
}

func TestMemoryWriteCanGrow(t *testing.T) {
	NewConformanceWithPath(setupMemory(platform.Linux)).TestWrite(t,
		"/gran/parent/child",
		"grow.txt",
		[]byte("this is test data"),
		7,
		[]byte(" more data than expected"),
		[]byte("this is more data than expected"))
}

func TestMemoryWriteCanOverwriteMiddle(t *testing.T) {
	NewConformanceWithPath(setupMemory(platform.Linux)).TestWrite(t,
		"/gran/parent/child",
		"less.txt",
		[]byte("this is test data"),
		8,
		[]byte("also"),
		[]byte("this is also data"))
}

func TestMemoryWriteCanOverwriteEnd(t *testing.T) {
	NewConformanceWithPath(setupMemory(platform.Linux)).TestWrite(t,
		"/gran/parent/child",
		"end.txt",
		[]byte("this is test data"),
		13,
		[]byte("info"),
		[]byte("this is test info"))

}

func TestMemoryReadDir(t *testing.T) {
	NewConformanceWithPath(setupMemory(platform.Linux)).
		TestReadDir(t,
			"/gran/parent/child", []file{
				{"one.txt", []byte("one")},
				{"two.txt", []byte("two")},
				{"three.txt", []byte("three")},
			})
}

func TestMemoryCanCreateFile(t *testing.T) {
	NewConformanceWithPath(setupMemory(platform.Linux)).
		TestCanCreateFile(t, "/gran/parent/child", []file{
			{"one.txt", []byte("one")},
			{"two.txt", []byte("two")},
			{"three.txt", []byte("three")},
		})
}

func TestMemoryCanWriteFile(t *testing.T) {
	NewConformanceWithPath(setupMemory(platform.Linux)).
		TestCanWriteFile(t, "/gran/parent/child", []file{
			{"one.txt", []byte("one")},
			{"two.txt", []byte("two")},
			{"three.txt", []byte("three")},
		})
}

func TestMemoryOpenFileFailsWhenReadOnlyAndNotExists(t *testing.T) {
	NewConformanceWithPath(setupMemory(platform.Linux)).
		TestOpenFileFailsWhenNotExists(t, "/gran/parent/child", "/gran/parent/child/one.txt")

}

func TestWindowsWillNormalizePath(t *testing.T) {
	NewConformanceWithPath(setupMemory(platform.Windows)).
		TestWindowsWillNormalizePath(t, `c:/ProgramData/fake/folder`, `test.txt`)
}

func TestWindowsFileExists(t *testing.T) {
	NewConformanceWithPath(setupMemory(platform.Windows)).
		TestWindowsFileForwardAndBackwardSlash(t, "c:/ProgramData/fake/folder/test.txt")
}

func setupMemory(plat platform.Platform) (fs.FS, filepath.Processor) {
	processor := filepath.NewProcessorWithPlatform(plat)
	fs := fs.NewMemory(fs.WithProcessor(processor))
	return fs, processor
}
