package xfilepath

import (
	"os"
	"path/filepath"
	"strings"
)

type JoinPath interface {
	Join(elements ...string) string
}

type RelPath interface {
	Rel(sourcepath, targetpath string) (string, error)
}

type CleanPath interface {
	Clean(path string) string
}

type RootPath interface {
	Root(path string) string
}

type VolumeNamePath interface {
	VolumeName(path string) string
}

type Processor interface {
	JoinPath
	RelPath
	CleanPath
	RootPath
	VolumeNamePath
	Separator() PathSeparator
}

type processor struct {
	sep PathSeparator
}

func NewDefaultProcessor() Processor {
	return NewProcessorWith(Default)
}

func NewProcessorWith(sep PathSeparator) Processor {
	return &processor{
		sep: sep,
	}
}

// Join implements Processor
func (p *processor) Join(elements ...string) string {
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

	return p.String(accumulator)
}

// Rel implements Processor
func (p *processor) Rel(sourcepath string, targetpath string) (string, error) {
	rel, err := filepath.Rel(sourcepath, targetpath)
	if err != nil {
		return "", err
	}
	return p.Clean(rel), nil
}

// Clean implements Processor
func (p *processor) Clean(path string) string {
	clean := filepath.Clean(path)
	return strings.ReplaceAll(clean, string(os.PathSeparator), string(p.sep))
}

// Root is a helper function to print the root of the filepath
func (p *processor) Root(path string) string {
	fp, _ := Parse(path)
	return p.String(fp.Root())
}

// VolumeName behaves similar to filepath.VolumeName in the path/filepath package
func (p *processor) VolumeName(path string) string {
	fp, _ := Parse(path)
	return fp.VolumeName(p.sep)
}

func (p *processor) String(fp FilePath) string {
	return fp.String(p.sep)
}

func (p *processor) Separator() PathSeparator {
	return p.sep
}
