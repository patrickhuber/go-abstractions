package os

import (
	"runtime"

	"github.com/patrickhuber/go-xplat/arch"
	"github.com/patrickhuber/go-xplat/platform"
)

const (
	MockAmd64Architecture = arch.AMD64
	MockArm64Architecture = arch.ARM64be

	MockWindowsPlatform         = platform.Windows
	MockWindowsWorkingDirectory = "c:\\working"
	MockWindowsHomeDirectory    = "c:\\users\\fake"

	MockLinuxPlatform        = platform.Linux
	MockUnixWorkingDirectory = "/working"
	MockUnixHomeDirectory    = "/home/fake"

	MockDarwinPlatform = platform.Darwin
)

type mockOS struct {
	workingDirectory string
	platform         platform.Platform
	architecture     arch.Arch
	homeDirectory    string
}

type MockOption func(*mockOS)

func WithHomeDirectory(homeDirectory string) MockOption {
	return func(o *mockOS) {
		o.homeDirectory = homeDirectory
	}
}

func WithArchitecture(architecture arch.Arch) MockOption {
	return func(o *mockOS) {
		o.architecture = architecture
	}
}

func WithWorkingDirectory(workingDirectory string) MockOption {
	return func(o *mockOS) {
		o.workingDirectory = workingDirectory
	}
}

func WithPlatform(platform platform.Platform) MockOption {
	return func(o *mockOS) {
		o.platform = platform
	}
}

// NewMock creates a new OS from the mock OS request
func NewMock(options ...MockOption) OS {
	o := &mockOS{}
	for _, option := range options {
		option(o)
	}
	if o.architecture == arch.Arch("") {
		o.architecture = arch.Arch(runtime.GOARCH)
	}
	if o.platform == platform.Platform("") {
		o.platform = platform.Platform(runtime.GOOS)
	}
	if o.workingDirectory == "" {
		if o.platform.IsWindows() {
			o.workingDirectory = MockWindowsWorkingDirectory
		} else {
			o.workingDirectory = MockUnixWorkingDirectory
		}
	}
	if o.homeDirectory == "" {
		if o.platform.IsWindows() {
			o.homeDirectory = MockWindowsHomeDirectory
		} else {
			o.homeDirectory = MockUnixHomeDirectory
		}
	}
	return o
}

func (o *mockOS) WorkingDirectory() (string, error) {
	return o.workingDirectory, nil
}

func (o *mockOS) Platform() platform.Platform {
	return platform.Platform(o.platform)
}

func (o *mockOS) Architecture() arch.Arch {
	return o.architecture
}

func (o *mockOS) Home() string {
	return o.homeDirectory
}

func (o *mockOS) ChangeDirectory(dir string) error {
	o.workingDirectory = dir
	return nil
}
