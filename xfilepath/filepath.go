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
		return fp.uncVolumeName(sep)

	case fp.IsRel():
		// Relative
		return ""
	}

	// Unix
	return ""
}

func (fp FilePath) uncVolumeName(sep PathSeparator) string {
	var builder strings.Builder

	// write //
	builder.WriteByte(byte(sep))
	builder.WriteByte(byte(sep))

	// write the hostname
	if fp.Volume.Host.HasValue {
		builder.WriteString(fp.Volume.Host.Value)
	}

	// write separator and share if share exists
	if fp.Volume.Share.HasValue {
		builder.WriteByte(byte(sep))
		builder.WriteString(fp.Volume.Share.Value)
	}

	return builder.String()
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

	// unc paths with one empty segment are already clean
	if fp.isUNC() && len(fp.Segments) == 1 && fp.Segments[0] == "" {
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

	return FilePath{
		Volume:   fp.Volume,
		Absolute: fp.Absolute,
		Segments: toSlice(s),
	}
}

func toSlice[T any](s stack.Stack[T]) []T {
	var slice []T
	for {
		if s.Length() == 0 {
			break
		}
		item := s.Pop()
		slice = append([]T{item}, slice...)
	}
	return slice
}

func (fp FilePath) Rel(other FilePath) (FilePath, error) {
	source := fp.Clean()
	target := other.Clean()

	// if paths are equal, return CurrentDirectory string
	if source.Equal(target) {
		return FilePath{
			Segments: []string{CurrentDirectory},
		}, nil
	}

	// remove any current directory '.' only source paths
	if len(source.Segments) == 1 && source.Segments[0] == CurrentDirectory {
		source.Segments = nil
	}

	// both paths must be either relative or absolute
	// if absolute both paths must match volumes
	if source.Absolute != target.Absolute || !source.Volume.Equal(target.Volume) {
		return FilePath{}, fmt.Errorf("can't make target relative to source: absolute paths must share a prefix")
	}

	// get the first index where segments differ
	firstDiff := source.firstSegmentDiff(target)

	if firstDiff < len(source.Segments) && source.Segments[firstDiff] == ".." {
		return FilePath{}, fmt.Errorf("can't make target relative to source")
	}

	var segments []string

	// run the source to the end by adding ..
	for s := firstDiff; s < len(source.Segments); s++ {
		segments = append(segments, ParentDirectory)
	}

	// run the target to the end by adding target[firstDiff:]
	if firstDiff < len(target.Segments) {
		segments = append(segments, target.Segments[firstDiff:]...)
	}

	return FilePath{
		Segments: segments,
		Absolute: false,
	}, nil
}

func (source FilePath) firstSegmentDiff(target FilePath) int {

	sourceLen := len(source.Segments)
	targetLen := len(target.Segments)
	segmentLen := sourceLen
	if targetLen < sourceLen {
		segmentLen = targetLen
	}

	diffPosition := 0

	// find the first differing element
	for ; diffPosition < segmentLen; diffPosition++ {
		if source.Segments[diffPosition] != target.Segments[diffPosition] {
			break
		}
	}
	return diffPosition
}

func (fp FilePath) Equal(other FilePath) bool {
	if fp.Absolute != other.Absolute {
		return false
	}
	if !fp.Volume.Equal(other.Volume) {
		return false
	}
	if len(fp.Segments) != len(other.Segments) {
		return false
	}
	for i := 0; i < len(fp.Segments); i++ {
		if fp.Segments[i] != other.Segments[i] {
			return false
		}
	}
	return true
}

func (v Volume) Equal(other Volume) bool {
	if !v.Drive.Equal(other.Drive) {
		return false
	}
	if !v.Host.Equal(other.Host) {
		return false
	}
	return v.Share.Equal(other.Share)
}

func (s NullableString) Equal(other NullableString) bool {
	return s.HasValue == other.HasValue &&
		s.Value == other.Value
}
