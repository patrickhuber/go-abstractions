package xfilepath

import (
	"os"
	"strings"
)

type PathType string

type FilePath struct {
	Volume   Volume
	Absolute bool
	Segments []string
}

type Volume struct {
	Host  string
	Share string
	Drive string
}

type PathSeparator rune

const (
	ForwardSlash  PathSeparator = '/'
	BackwardSlash PathSeparator = '\\'
	Default       PathSeparator = os.PathSeparator
)

func (fp FilePath) Root(sep PathSeparator) string {

	switch {
	case !fp.Absolute:
		// Relative
		return ""

	case fp.Volume.Drive != "":
		// Windows
		var builder strings.Builder
		builder.WriteString(fp.Volume.Drive)
		return builder.String()

	case fp.Volume.Host != "":
		// UNC
		var builder strings.Builder
		builder.WriteByte(byte(sep))
		builder.WriteByte(byte(sep))
		builder.WriteString(fp.Volume.Host)
		builder.WriteByte(byte(sep))
		builder.WriteString(fp.Volume.Share)
		return builder.String()
	}

	// Unix
	var builder strings.Builder
	builder.WriteByte(byte(sep))
	return builder.String()
}

// VolumeName behaves similar to filepath.VolumeName in the path/filepath package
func (fp FilePath) VolumeName(sep PathSeparator) string {
	switch {
	case !fp.Absolute:
		// Relative
		return ""

	case fp.Volume.Drive != "":
		// Windows
		var builder strings.Builder
		builder.WriteString(fp.Volume.Drive)
		return builder.String()

	case fp.Volume.Host != "":
		// UNC
		var builder strings.Builder
		builder.WriteByte(byte(sep))
		builder.WriteByte(byte(sep))
		builder.WriteString(fp.Volume.Host)
		builder.WriteByte(byte(sep))
		builder.WriteString(fp.Volume.Share)
		return builder.String()
	}

	// Unix
	return ""
}

func (fp FilePath) String(sep PathSeparator) string {

	volumeName := fp.VolumeName(sep)

	var builder strings.Builder
	builder.WriteString(volumeName)

	for i, seg := range fp.Segments {
		// print the prefix slash if the
		// path is not relative
		// or this is not the first element
		if i != 0 || fp.Absolute {
			builder.WriteByte(byte(sep))
		}
		builder.WriteString(seg)
	}

	return builder.String()
}

func (fp FilePath) Join(other FilePath) FilePath {
	return FilePath{
		Volume: Volume{
			Host:  fp.Volume.Host,
			Share: fp.Volume.Share,
			Drive: fp.Volume.Drive,
		},
		Absolute: fp.Absolute,
		Segments: append(fp.Segments, other.Segments...),
	}
}

// Join is a helper function to combine elements into a single FilePath object
// it works by parsing each element and sequentially appending the elements together
// it then calls String on the result
func Join(sep PathSeparator, elements ...string) string {

	if len(elements) == 0 {
		return ""
	}

	var accumulator FilePath
	first := true
	for _, element := range elements {

		// skip empty elements
		if len(element) == 0 {
			continue
		}

		// call parse on the first element
		// set the first element as the accumulator
		if first {
			accumulator, _ = Parse(element)
			first = false
			continue
		}

		// call parse on each next element
		next, _ := Parse(element)

		// and then join the accumulator to that element
		accumulator = accumulator.Join(next)
	}

	return accumulator.String(sep)
}
