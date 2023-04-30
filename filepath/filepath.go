package filepath

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
type Comparison int

const (
	ForwardSlash             PathSeparator     = '/'
	BackwardSlash            PathSeparator     = '\\'
	Colon                    PathListSeparator = ':'
	SemiColon                PathListSeparator = ';'
	DefaultPathListSeparator PathListSeparator = os.PathListSeparator
	DefaultPathSeparator     PathSeparator     = os.PathSeparator
	IgnoreCase               Comparison        = 1
	CaseSensitive            Comparison        = 0

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

func (fp FilePath) Rel(other FilePath, cmp Comparison) (FilePath, error) {
	source := fp.Clean()
	target := other.Clean()

	// if paths are equal, return CurrentDirectory string
	if source.Equal(target, cmp) {
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
	if source.Absolute != target.Absolute || !source.Volume.Equal(target.Volume, cmp) {
		return FilePath{}, fmt.Errorf("can't make target relative to source: absolute paths must share a prefix")
	}

	// get the first index where segments differ
	firstDiff := source.firstSegmentDiff(target, cmp)

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

func (source FilePath) firstSegmentDiff(target FilePath, cmp Comparison) int {

	sourceLen := len(source.Segments)
	targetLen := len(target.Segments)

	segmentLen := sourceLen
	if targetLen < sourceLen {
		segmentLen = targetLen
	}

	// find the first differing element
	for diffPosition := 0; diffPosition < segmentLen; diffPosition++ {
		if !equal(source.Segments[diffPosition], target.Segments[diffPosition], cmp) {
			return diffPosition
		}
	}
	return segmentLen
}

func (fp FilePath) Dir() FilePath {
	last := len(fp.Segments) - 1
	if last >= 0 && fp.Segments[last] != "" {
		fp.Segments[last] = ""
	}
	return fp.Clean()
}

func (fp FilePath) Ext() string {
	if len(fp.Segments) == 0 {
		return ""
	}
	last := fp.Segments[len(fp.Segments)-1]
	for i := len(last) - 1; i >= 0; i-- {
		if last[i] == '.' {
			return last[i:]
		}
	}
	return ""
}

// Equal compares two paths using case sensitive comparison
func (fp FilePath) Equal(other FilePath, cmp Comparison) bool {
	if fp.Absolute != other.Absolute {
		return false
	}
	if !fp.Volume.Equal(other.Volume, cmp) {
		return false
	}
	if len(fp.Segments) != len(other.Segments) {
		return false
	}

	for i := 0; i < len(fp.Segments); i++ {
		if !equal(fp.Segments[i], other.Segments[i], cmp) {
			return false
		}
	}
	return true
}

func equal(s, t string, cmp Comparison) bool {
	if cmp == IgnoreCase {
		return strings.EqualFold(s, t)
	}
	return s == t
}

// Equal compares two volumes using case sensetive comparison
func (v Volume) Equal(other Volume, cmp Comparison) bool {
	if !v.Drive.Equal(other.Drive, cmp) {
		return false
	}
	if !v.Host.Equal(other.Host, cmp) {
		return false
	}
	return v.Share.Equal(other.Share, cmp)
}

// Equal compares two nullable strings using case sensitive comparison
func (s NullableString) Equal(other NullableString, cmp Comparison) bool {
	// both strings must have a value or not have a value
	if s.HasValue != other.HasValue {
		return false
	}

	// both do not have value
	if !s.HasValue {
		return true
	}

	// both have value, return the equality of the values
	return equal(s.Value, other.Value, cmp)
}
