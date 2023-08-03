package filepath

import (
	"regexp"

	"github.com/patrickhuber/go-xplat/os"
	"github.com/patrickhuber/go-xplat/platform"
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

type DirPath interface {
	Dir(path string) string
}

type ExtPath interface {
	Ext(path string) string
}

type BasePath interface {
	Base(path string) string
}

type AbsPath interface {
	Abs(path string) (string, error)
}

type Processor interface {
	JoinPath
	RelPath
	CleanPath
	RootPath
	VolumeNamePath
	DirPath
	ExtPath
	BasePath
	AbsPath
	Separator() PathSeparator
	Parser() Parser
	Comparison() Comparison
}

type processor struct {
	platform platform.Platform
	os       os.OS
	sep      PathSeparator
	parser   Parser
	cmp      Comparison
}

type ProcessorOption func(p *processor)

func WithParser(parser Parser) ProcessorOption {
	return func(p *processor) {
		p.parser = parser
	}
}

func WithSeparator(separator PathSeparator) ProcessorOption {
	return func(p *processor) {
		p.sep = separator
	}
}

func WithComparison(cmp Comparison) ProcessorOption {
	return func(p *processor) {
		p.cmp = cmp
	}
}

func WithOS(os os.OS) ProcessorOption {
	return func(p *processor) {
		p.os = os
	}
}

// NewProcessor creates a processor with the default platform and then applies the options
func NewProcessor(options ...ProcessorOption) Processor {
	return NewProcessorWithPlatform(platform.Default(), options...)
}

// NewProcessorWithPlatform creates a platform specific processor and then applies the given options
func NewProcessorWithPlatform(plat platform.Platform, options ...ProcessorOption) Processor {

	p := &processor{
		parser:   NewParserWithPlatform(plat),
		platform: plat,
	}

	if plat.IsUnix() {
		p.sep = ForwardSlash
		p.cmp = CaseSensitive
	} else {
		p.sep = BackwardSlash
		p.cmp = IgnoreCase
	}

	for _, option := range options {
		option(p)
	}

	// set the default OS so we can use the working directory wrapper
	if p.os == nil {
		p.os = os.New()
	}

	return p
}

func (p *processor) Abs(path string) (string, error) {
	wd, err := p.os.WorkingDirectory()
	if err != nil {
		return "", err
	}
	return p.abs(wd, path)
}

func (p *processor) abs(wd, rel string) (string, error) {
	fp, err := p.parser.Parse(rel)
	if err != nil {
		return "", err
	}
	if fp.IsAbs() {
		return p.String(fp.Clean()), nil
	}
	wdp, err := p.parser.Parse(wd)
	if err != nil {
		return "", err
	}
	abs := wdp.Join(fp)
	return p.String(abs.Clean()), nil
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
			accumulator, _ = p.parser.Parse(element)
			first = false
			continue
		}

		// call parse on each next element
		next, _ := p.parser.Parse(element)

		// and then join the accumulator to that element
		accumulator = accumulator.Join(next)
	}

	// call clean on the result
	return p.String(accumulator.Clean())
}

// Rel implements Processor
func (p *processor) Rel(sourcepath string, targetpath string) (string, error) {

	source, err := p.parser.Parse(sourcepath)
	if err != nil {
		return "", err
	}

	target, err := p.parser.Parse(targetpath)
	if err != nil {
		return "", err
	}

	result, err := source.Rel(target, p.cmp)
	if err != nil {
		return "", err
	}

	return result.String(p.sep), nil
}

// Clean implements Processor
func (p *processor) Clean(path string) string {
	fp, _ := p.parser.Parse(path)

	// for empty unc paths, normalize the original string
	// (is there a way to do this in the String method?)
	fp = fp.Clean()

	// on the windows platform if the first segment matches the drive pattern
	// the current directory needs to be added in the front
	if p.platform == platform.Windows && fp.IsRel() && len(fp.Segments) > 0 {
		matched, err := regexp.MatchString(`^[a-zA-Z][:]`, fp.Segments[0])
		if (err == nil) && matched {
			fp.Segments = append([]string{CurrentDirectory}, fp.Segments...)
		}
	}

	cleaned := fp.String(p.sep)
	return cleaned
}

// Root is a helper function to print the root of the filepath
func (p *processor) Root(path string) string {
	fp, _ := p.parser.Parse(path)
	return p.String(fp.Root())
}

// VolumeName behaves similar to filepath.VolumeName in the path/filepath package
func (p *processor) VolumeName(path string) string {
	fp, _ := p.parser.Parse(path)
	return fp.VolumeName(p.sep)
}

func (p *processor) Ext(path string) string {
	fp, _ := p.parser.Parse(path)
	return fp.Ext()
}

func (p *processor) Dir(path string) string {
	fp, _ := p.parser.Parse(path)
	dir := fp.Dir()
	return dir.String(p.sep)
}

// Base returns the last element of path. Trailing path separators are removed before extracting the last element. If the path is empty, Base returns ".". If the path consists entirely of separators, Base returns a single separator.
func (p *processor) Base(path string) string {
	fp, _ := p.parser.Parse(path)
	base := fp.Base()
	return base.String(p.sep)
}

// String returns the string representation of the file path
func (p *processor) String(fp FilePath) string {
	return fp.String(p.sep)
}

// Separator returns the PathSeparator for the processor
func (p *processor) Separator() PathSeparator {
	return p.sep
}

// Parser returns the Parser for the processor
func (p *processor) Parser() Parser {
	return p.parser
}

// Comparison returns the path comparison operator for the processor
func (p *processor) Comparison() Comparison {
	return p.cmp
}
