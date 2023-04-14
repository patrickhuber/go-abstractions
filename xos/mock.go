package os

const (
	MockAmd64Architecture = "amd64"
	MockArm64Architecture = "arm64"

	MockWindowsPlatform         = "windows"
	MockWindowsWorkingDirectory = "c:\\working"
	MockWindowsHomeDirectory    = "c:\\users\\fake"
	MockWindowsExecutable       = "c:\\ProgramData\\test\\fake.exe"

	MockLinuxPlatform         = "linux"
	MockLinuxWorkingDirectory = "/working"
	MockLinuxHomeDirectory    = "/home/fake"
	MockLinuxExecutable       = "/opt/test/fake"

	MockDarwinPlatform         = "darwin"
	MockDarwinHomeDirectory    = MockLinuxHomeDirectory
	MockDarwinWorkingDirectory = MockLinuxWorkingDirectory
	MockDarwinExecutable       = MockLinuxExecutable
)

type mockOS struct {
	workingDirectory string
	platform         string
	architecture     string
	homeDirectory    string
	executable       string
}

type NewMockOS struct {
	WorkingDirectory string
	Platform         string
	Architecture     string
	HomeDirectory    string
	Executable       string
}

type MockOption func(*NewMockOS)

func WithExecutable(exectuable string) MockOption {
	return func(o *NewMockOS) {
		o.Executable = exectuable
	}
}

func WithHomeDirectory(homeDirectory string) MockOption {
	return func(o *NewMockOS) {
		o.HomeDirectory = homeDirectory
	}
}

func WithArchitecture(architecture string) MockOption {
	return func(o *NewMockOS) {
		o.Architecture = architecture
	}
}

func WithWorkingDirectory(workingDirectory string) MockOption {
	return func(o *NewMockOS) {
		o.WorkingDirectory = workingDirectory
	}
}

// NewMock creates a new OS from the mock OS request
func NewMock(o *NewMockOS, options ...MockOption) OS {
	return &mockOS{
		executable:       o.Executable,
		workingDirectory: o.WorkingDirectory,
		architecture:     o.Architecture,
		platform:         o.Platform,
		homeDirectory:    o.HomeDirectory,
	}
}

func NewLinuxMock(options ...MockOption) OS {
	mock := &NewMockOS{
		Executable:       MockLinuxExecutable,
		WorkingDirectory: MockLinuxWorkingDirectory,
		Platform:         MockLinuxPlatform,
		HomeDirectory:    MockLinuxHomeDirectory,
		Architecture:     MockAmd64Architecture,
	}
	return NewMock(mock, options...)
}

func NewDarwinMock(options ...MockOption) OS {
	mock := &NewMockOS{
		Executable:       MockDarwinExecutable,
		WorkingDirectory: MockDarwinWorkingDirectory,
		Platform:         MockDarwinPlatform,
		HomeDirectory:    MockDarwinHomeDirectory,
		Architecture:     MockAmd64Architecture,
	}
	return NewMock(mock, options...)
}

func NewWindowsMock(options ...MockOption) OS {
	mock := &NewMockOS{
		Executable:       MockWindowsExecutable,
		WorkingDirectory: MockWindowsWorkingDirectory,
		Platform:         MockWindowsPlatform,
		HomeDirectory:    MockWindowsHomeDirectory,
		Architecture:     MockAmd64Architecture,
	}
	return NewMock(mock, options...)
}

func (o *mockOS) WorkingDirectory() (string, error) {
	return o.workingDirectory, nil
}

func (o *mockOS) Executable() (string, error) {
	return o.executable, nil
}

func (o *mockOS) Platform() string {
	return o.platform
}

func (o *mockOS) Architecture() string {
	return o.architecture
}

func (o *mockOS) Home() string {
	return o.homeDirectory
}
