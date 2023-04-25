package xfilepath

import (
	"fmt"
	"os"
	"strings"

	"github.com/patrickhuber/go-collections/generic/stack"
)

type PathType string

type FilePath struct {
	Volume   Volume
	Absolute bool
	Segments []string
}

type Volume struct {
	Host  NullableString
	Share NullableString
	Drive NullableString
}

type NullableString struct {
	HasValue bool
	Value    string
}

type PathSeparator rune
type PathListSeparator rune

const (
	ForwardSlash             PathSeparator     = '/'
	BackwardSlash            PathSeparator     = '\\'
	Colon                    PathListSeparator = ':'
	SemiColon                PathListSeparator = ';'
	DefaultPathListSeparator PathListSeparator = os.PathListSeparator
	DefaultPathSeparator     PathSeparator     = os.PathSeparator

	CurrentDirectory = "."
	ParentDirectory  = ".."
	EmptyDirectory   = ""
)

func (fp FilePath) IsAbs() bool {
	return fp.Absolute
}

func (fp FilePath) IsRel() bool {
	return !fp.Absolute
}

func (fp FilePath) isWindows() bool {
	return fp.Volume.Drive.HasValue
}

func (fp FilePath) isUNC() bool {
	return fp.Volume.Host.HasValue
}

func (fp FilePath) Root() FilePath {
	return FilePath{
		Volume:   fp.Volume,
		Absolute: fp.Absolute,
	}
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

func (fp FilePath) VolumeName(sep PathSeparator) string {
	switch {

	case fp.isWindows():
		// Windows
		var builder strings.Builder
		builder.WriteString(fp.Volume.Drive.Value)
		return builder.String()

	case fp.isUNC():
		// UNC
		var builder strings.Builder
		builder.WriteByte(byte(sep))
		builder.WriteByte(byte(sep))
		if fp.Volume.Host.HasValue {
			builder.WriteString(fp.Volume.Host.Value)
			builder.WriteByte(byte(sep))
		}
		if fp.Volume.Share.HasValue {
			builder.WriteString(fp.Volume.Share.Value)
		}
		return builder.String()

	case fp.IsRel():
		// Relative
		return ""
	}

	// Unix
	return ""
}

func (fp FilePath) String(sep PathSeparator) string {
	var builder strings.Builder

	// write the volume name
	builder.WriteString(fp.VolumeName(sep))

	switch {
	// relative paths don't need a separator
	case fp.IsRel():

	// absolute windows and unix paths need a separator
	case !fp.isUNC():
		builder.WriteRune(rune(sep))

	// unc paths with segments need a separator
	case len(fp.Segments) > 0:
		builder.WriteRune(rune(sep))
	}

	// write the segments
	for i, seg := range fp.Segments {
		if i > 0 {
			builder.WriteRune(rune(sep))
		}
		builder.WriteString(seg)
	}
	return builder.String()
}

func (fp FilePath) Clean() FilePath {
	clean := fp.clean()

	// add a current directory indicator for empty relative paths
	if len(clean.Segments) == 0 && clean.IsRel() {
		clean.Segments = append(clean.Segments, CurrentDirectory)
	}
	return clean
}

func (fp FilePath) clean() FilePath {
	// if the path has no segments it is already clean
	if len(fp.Segments) == 0 {
		return fp
	}

	s := stack.New[string]()
	for _, segment := range fp.Segments {

		switch segment {
		// remove . and empty directories (resulting from // in path)
		case CurrentDirectory:
			continue
		case EmptyDirectory:
			continue
		case ParentDirectory:
			switch {
			case s.Length() > 0:
				// if the current segment is parent
				// and the last segment is parent
				// we have already processed a parent so write both to the output
				if s.Pop() == ParentDirectory {
					s.Push(ParentDirectory)
					s.Push(ParentDirectory)
				}
			case fp.IsRel():
				// if the current segment is parent
				// and there are no elements to process (s.Length() == 0)
				// push the ..
				s.Push(ParentDirectory)
			}
		default:
			// if the segment matches the pattern \w[:] it is a drive letter in windows
			s.Push(segment)
		}
	}
	var segments []string
	for {
		if s.Length() == 0 {
			break
		}
		item := s.Pop()
		segments = append([]string{item}, segments...)
	}

	return FilePath{
		Volume:   fp.Volume,
		Absolute: fp.Absolute,
		Segments: segments,
	}
}

func (fp FilePath) Rel(other FilePath) (FilePath, error) {
	return FilePath{}, fmt.Errorf("not implemented")
}
