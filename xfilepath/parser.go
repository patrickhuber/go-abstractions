package xfilepath

import (
	"strings"
	"unicode"

	"github.com/patrickhuber/go-xplat/platform"
)

type PathListParser interface {
	ListSeparator() PathListSeparator
	ParseList(path string) ([]FilePath, error)
}

type PathParser interface {
	Separators() []PathSeparator
	Parse(path string) (FilePath, error)
}

type Parser interface {
	PathParser
	PathListParser
}

type ParseOption func(*parser)

func ifEmptyReturnNil(slice []string) []string {
	if len(slice) == 0 {
		return nil
	}
	return slice
}

type parser struct {
	platform      platform.Platform
	listSeparator PathListSeparator
	separators    []PathSeparator
}

// NewParser creates a parser for the current platform and then applies options
func NewParser(options ...ParseOption) Parser {
	// set the default platform
	return NewParserWithPlatform(platform.Default(), options...)
}

// NewParserWithPlatform creates a new parser with the specified platform defaults and then applies options
func NewParserWithPlatform(plat platform.Platform, options ...ParseOption) Parser {
	p := &parser{
		platform: plat,
	}
	if plat.IsUnix() {
		p.separators = []PathSeparator{ForwardSlash}
		p.listSeparator = Colon
	} else {
		p.separators = []PathSeparator{BackwardSlash, ForwardSlash}
		p.listSeparator = SemiColon
	}
	for _, option := range options {
		option(p)
	}
	return p
}

func WithListSeparator(sep PathListSeparator) ParseOption {
	return func(p *parser) {
		p.listSeparator = sep
	}
}

func WithSeparators(sep ...PathSeparator) ParseOption {
	return func(p *parser) {
		p.separators = sep
	}
}

func (p *parser) Separators() []PathSeparator {
	return p.separators
}

func (p *parser) ListSeparator() PathListSeparator {
	return p.listSeparator
}

func (p *parser) Parse(path string) (FilePath, error) {
	// if it is a UNC path
	if p.isUNCPath(path) {
		return p.parseUNCPath(path)
	}
	if p.isWindowsPath(path) {
		return p.parseWindowsPath(path)
	}
	return p.parseUnixPath(path)
}

func (p *parser) ParseList(path string) ([]FilePath, error) {
	var list []FilePath
	for {
		before, after, found := strings.Cut(path, string(p.listSeparator))

		fp, err := p.Parse(before)
		if err != nil {
			return nil, err
		}
		list = append(list, fp)
		if !found {
			break
		}
		path = after
	}
	return list, nil
}

func (p *parser) parseUNCPath(path string) (FilePath, error) {
	// skip the two slashes
	path = path[2:]
	segments := p.split(path)

	var host NullableString
	if len(segments) > 0 {
		host = NullableString{
			Value:    segments[0],
			HasValue: true,
		}
	}

	var share NullableString
	if len(segments) > 1 {
		share = NullableString{
			Value:    segments[1],
			HasValue: true,
		}
	}

	if len(segments) > 2 {
		segments = segments[2:]
	} else {
		segments = nil
	}

	return FilePath{
		Volume: Volume{
			Host:  host,
			Share: share,
		},
		Absolute: true,
		Segments: segments}, nil
}

func (p *parser) parseWindowsPath(path string) (FilePath, error) {
	// remove the drive letter from the path and get the path segments
	segments := p.split(path[2:])

	// remove the preceeding empty space
	if len(segments) > 1 {
		if segments[0] == "" {
			segments = segments[1:]
		}
	}

	return FilePath{
		Volume: Volume{
			Drive: NullableString{Value: path[0:2], HasValue: true},
		},
		Absolute: len(path) > 2 && p.isSeparator(path[2]),
		Segments: ifEmptyReturnNil(segments),
	}, nil
}

func (p *parser) parseUnixPath(path string) (FilePath, error) {
	segments := p.split(path)
	absolute := false
	if len(path) > 0 {
		absolute = p.isSeparator(path[0])
	}

	// remove the preceeding empty space in absoluate paths
	if len(segments) > 1 && absolute {
		if segments[0] == "" {
			segments = segments[1:]
		}
	}

	// special case for "/"
	if len(segments) == 1 && absolute && segments[0] == "" {
		segments = nil
	}

	return FilePath{
		Segments: ifEmptyReturnNil(segments),
		Absolute: absolute,
	}, nil
}

func (p *parser) isWindowsPath(path string) bool {
	// or the drive letter and colon exist
	return p.isDrive(path)
}

func (p *parser) isUNCPath(path string) bool {
	// the first two slashes could be a unix path with an empty path element
	if len(path) <= 2 {
		return false
	}
	if p.platform != platform.Windows {
		return false
	}
	return p.isSeparator(path[0]) && p.isSeparator(path[1])
}

func (p *parser) isDrive(path string) bool {
	if len(path) < 2 {
		return false
	}
	if !unicode.IsLetter(rune(path[0])) {
		return false
	}
	return path[1] == ':'
}

// Split splits a path into segments preserving leading and trailing empty segments
//
// given "" returns nil
// given "/" returns ["", ""]
// given "something/" returns ["something", ""]
// given "/something" returns ["", "something"]
// given "something" returns ["something"]
// given "/something/" returns ["", "something", ""]
func (p *parser) split(path string) []string {

	// empty string
	var segments []string
	if len(path) == 0 {
		return segments
	}

	// split will contain the list of segments and empty segments where two separators are adjacent
	for {
		before, ok, after := p.cut(path)
		segments = append(segments, before)
		if !ok {
			break
		}
		path = after
	}
	return segments
}

// cut cuts the path at the first separator occurence
// before contains the string before the first separator
// found returns true if a separator was found, false otherwise
// after contains the string after the first separator
func (p *parser) cut(path string) (before string, found bool, after string) {
	for i := 0; i < len(path); i++ {
		if p.isSeparator(path[i]) {
			return path[0:i], true, path[i+1:]
		}
	}
	return path, false, ""
}

func (p *parser) isSeparator(b byte) bool {
	for _, sep := range p.separators {
		if b == byte(sep) {
			return true
		}
	}
	return false
}
