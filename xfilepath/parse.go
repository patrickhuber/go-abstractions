package xfilepath

import (
	"fmt"
	"unicode"
)

// FilePath = ('/'|'\){ 2 } Volume Segments
// | [a-zA-Z] ':' Segments
// | Segments
// Segments = '/' { Segment }
func Parse(path string) (FilePath, error) {
	// if it is a UNC path
	if isUNCPath(path) {
		return parseUNCPath(path)
	}
	if isWindowsPath(path) {
		return parseWindowsPath(path)
	}
	return parseUnixPath(path)
}

func parseUNCPath(path string) (FilePath, error) {
	originalPath := path
	// skip the two slashes
	path = path[2:]
	segments := split(path)
	if len(segments) < 2 {
		return FilePath{}, fmt.Errorf("invalid UNC path '%s'", originalPath)
	}

	return FilePath{
		Volume: Volume{
			Host:  segments[0],
			Share: segments[1],
		},
		Absolute: true,
		Trailing: hasTrailingSlash(path),
		Segments: ifEmptyReturnNil(segments[2:])}, nil
}

func hasTrailingSlash(path string) bool {
	if len(path) == 0 {
		return false
	}
	b := path[len(path)-1]
	return isSlash(b)
}

func parseWindowsPath(path string) (FilePath, error) {
	segments := split(path[2:])
	return FilePath{
		Volume: Volume{
			Drive: path[0:2],
		},
		Absolute: true,
		Trailing: hasTrailingSlash(path),
		Segments: ifEmptyReturnNil(segments),
	}, nil
}

func parseUnixPath(path string) (FilePath, error) {
	segments := split(path)
	absolute := true
	if len(path) > 0 {
		absolute = isSlash(path[0])
	}
	return FilePath{
		Segments: ifEmptyReturnNil(segments),
		Trailing: hasTrailingSlash(path),
		Absolute: absolute,
	}, nil
}

func isWindowsPath(path string) bool {

	// or the drive letter and colon exist
	return isDrive(path)
}

func isSlash(b byte) bool {
	return isForwardSlash(b) || isBackSlash(b)
}

func isBackSlash(b byte) bool {
	return b == byte(BackwardSlash)
}

func isForwardSlash(b byte) bool {
	return b == byte(ForwardSlash)
}

func isUNCPath(path string) bool {
	if len(path) <= 2 {
		return false
	}
	return isSlash(path[0]) && isSlash(path[1])
}

func isDrive(path string) bool {
	if len(path) < 2 {
		return false
	}
	if !unicode.IsLetter(rune(path[0])) {
		return false
	}
	return path[1] == ':'
}

// '/' returns {}
// '/something' returns {"something"}
// 'something/' returns {"something"}
// 'something' returns {"something"}
func split(path string) []string {
	var segments []string

	// empty string
	if len(path) == 0 {
		return segments
	}

	// remove leading slash
	if isSlash(path[0]) {
		if len(path) == 1 {
			return segments
		}
		path = path[1:]
	}

	for {
		before, ok, after := cut(path)

		// case "/"
		if before == "" && after == "" {
			return segments
		}

		segments = append(segments, before)
		if !ok {
			break
		}
		path = after
	}
	return segments
}

func cut(path string) (before string, found bool, after string) {
	for i := 0; i < len(path); i++ {
		if isSlash(path[i]) {
			return path[0:i], true, path[i+1:]
		}
	}
	return path, false, ""
}

func ifEmptyReturnNil(slice []string) []string {
	if len(slice) == 0 {
		return nil
	}
	return slice
}
