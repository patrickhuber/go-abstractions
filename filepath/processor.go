package filepath

import (
	"regexp"

	"github.com/patrickhuber/go-xplat/os"
)

type Processor struct {
	OS         os.OS
	Separator  PathSeparator
	Parser     Parser
	Comparison Comparison
}

type ProcessorOption func(p *Processor)

// NewProcessor creates a processor with the default platform and then applies the options
func NewProcessor() *Processor {
	return NewProcessorWithOS(os.New())
}

// NewProcessorWithOS creates a processor from the OS
func NewProcessorWithOS(o os.OS) *Processor {
	p := &Processor{
		Parser: NewParserWithPlatform(o.Platform()),
		OS:     o,
	}

	// run defaults after all options have passed
	if o.Platform().IsUnix() {
		p.Separator = ForwardSlash
		p.Comparison = CaseSensitive
	} else {
		p.Separator = BackwardSlash
		p.Comparison = IgnoreCase
	}
	return p
}

func (p *Processor) Abs(path string) (string, error) {
	wd, err := p.OS.WorkingDirectory()
	if err != nil {
		return "", err
	}
	return p.abs(wd, path)
}

func (p *Processor) abs(wd, rel string) (string, error) {
	fp, err := p.Parser.Parse(rel)
	if err != nil {
		return "", err
	}
	if fp.IsAbs() {
		return p.String(fp.Clean()), nil
	}
	wdp, err := p.Parser.Parse(wd)
	if err != nil {
		return "", err
	}
	abs := wdp.Join(fp)
	return p.String(abs.Clean()), nil
}

// Join implements Processor
func (p *Processor) Join(elements ...string) string {
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
			accumulator, _ = p.Parser.Parse(element)
			first = false
			continue
		}

		// call parse on each next element
		next, _ := p.Parser.Parse(element)

		// and then join the accumulator to that element
		accumulator = accumulator.Join(next)
	}

	// call clean on the result
	return p.String(accumulator.Clean())
}

// Rel implements Processor
func (p *Processor) Rel(sourcepath string, targetpath string) (string, error) {

	source, err := p.Parser.Parse(sourcepath)
	if err != nil {
		return "", err
	}

	target, err := p.Parser.Parse(targetpath)
	if err != nil {
		return "", err
	}

	result, err := source.Rel(target, p.Comparison)
	if err != nil {
		return "", err
	}

	return result.String(p.Separator), nil
}

// Clean implements Processor
func (p *Processor) Clean(path string) string {
	fp, _ := p.Parser.Parse(path)

	// for empty unc paths, normalize the original string
	// (is there a way to do this in the String method?)
	fp = fp.Clean()

	// on the windows platform if the first segment matches the drive pattern
	// the current directory needs to be added in the front
	if p.OS.Platform().IsWindows() && fp.IsRel() && len(fp.Segments) > 0 {
		matched, err := regexp.MatchString(`^[a-zA-Z][:]`, fp.Segments[0])
		if (err == nil) && matched {
			fp.Segments = append([]string{CurrentDirectory}, fp.Segments...)
		}
	}

	cleaned := fp.String(p.Separator)
	return cleaned
}

// Root is a helper function to print the root of the filepath
func (p *Processor) Root(path string) string {
	fp, _ := p.Parser.Parse(path)
	return p.String(fp.Root())
}

// VolumeName behaves similar to filepath.VolumeName in the path/filepath package
func (p *Processor) VolumeName(path string) string {
	fp, _ := p.Parser.Parse(path)
	return fp.VolumeName(p.Separator)
}

func (p *Processor) Ext(path string) string {
	fp, _ := p.Parser.Parse(path)
	return fp.Ext()
}

func (p *Processor) Dir(path string) string {
	fp, _ := p.Parser.Parse(path)
	dir := fp.Dir()
	return dir.String(p.Separator)
}

// Base returns the last element of path. Trailing path separators are removed before extracting the last element. If the path is empty, Base returns ".". If the path consists entirely of separators, Base returns a single separator.
func (p *Processor) Base(path string) string {
	fp, _ := p.Parser.Parse(path)
	base := fp.Base()
	return base.String(p.Separator)
}

// String returns the string representation of the file path
func (p *Processor) String(fp FilePath) string {
	return fp.String(p.Separator)
}
