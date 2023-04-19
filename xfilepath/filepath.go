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
	Trailing bool
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

func (fp FilePath) IsAbs() bool {
	return fp.Absolute
}

func (fp FilePath) IsRel() bool {
	return !fp.Absolute
}

func (fp FilePath) isWindows() bool {
	return fp.Volume.Drive != ""
}

func (fp FilePath) isUNC() bool {
	return fp.Volume.Host != ""
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
	case fp.IsRel():
		// Relative
		return ""

	case fp.isWindows():
		// Windows
		var builder strings.Builder
		builder.WriteString(fp.Volume.Drive)
		return builder.String()

	case fp.isUNC():
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
	var builder strings.Builder
	switch {

	case fp.IsRel():
		// we know that relative paths have no prefix

	case fp.isWindows() || fp.isUNC():
		builder.WriteString(fp.VolumeName(sep))
		if len(fp.Segments) > 0 {
			builder.WriteRune(rune(sep))
		}

	default:
		// unix
		builder.WriteRune(rune(sep))
	}

	for i, seg := range fp.Segments {
		if i > 0 {
			builder.WriteRune(rune(sep))
		}
		builder.WriteString(seg)
	}
	return builder.String()
}

func (fp FilePath) Clean() FilePath {
	// if the path has no segments
	if len(fp.Segments) == 0 {
		// if the path is a unc path
		return fp
	}
	for _, segment := range fp.Segments {
		switch segment {
		case ".":
		case "..":
		case "":
		}
	}
	return FilePath{}
}
