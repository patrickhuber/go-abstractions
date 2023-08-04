package os_test

import (
	"testing"

	"github.com/patrickhuber/go-xplat/arch"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/stretchr/testify/require"
)

func TestPlatform(t *testing.T) {
	type test struct {
		expected platform.Platform
		o        os.OS
	}
	tests := []test{
		{expected: os.MockDarwinPlatform, o: os.NewMock(os.WithPlatform(platform.Darwin))},
		{expected: os.MockLinuxPlatform, o: os.NewMock(os.WithPlatform(platform.Linux))},
		{expected: os.MockWindowsPlatform, o: os.NewMock(os.WithPlatform(platform.Windows))},
		{expected: platform.Plan9, o: os.NewMock(os.WithPlatform(platform.Plan9))},
		{expected: platform.Plan9, o: os.NewMock(os.WithPlatform(platform.Plan9))},
		{expected: platform.Plan9, o: os.NewMock(os.WithPlatform(platform.Plan9))},
		{expected: platform.Plan9, o: os.NewMock(os.WithPlatform(platform.Plan9))},
	}
	for _, test := range tests {
		require.Equal(t, test.expected, test.o.Platform())
	}
}

func TestArchitecture(t *testing.T) {
	type test struct {
		expected arch.Arch
		o        os.OS
	}
	tests := []test{
		{expected: os.MockAmd64Architecture, o: os.NewMock(os.WithPlatform(platform.Darwin))},
		{expected: os.MockAmd64Architecture, o: os.NewMock(os.WithPlatform(platform.Linux))},
		{expected: os.MockAmd64Architecture, o: os.NewMock(os.WithPlatform(platform.Windows))},
		{expected: os.MockArm64Architecture, o: os.NewMock(os.WithArchitecture(os.MockArm64Architecture))},
		{expected: os.MockArm64Architecture, o: os.NewMock(os.WithArchitecture(os.MockArm64Architecture))},
		{expected: os.MockArm64Architecture, o: os.NewMock(os.WithArchitecture(os.MockArm64Architecture))},
		{expected: os.MockArm64Architecture, o: os.NewMock(os.WithArchitecture(os.MockArm64Architecture))},
	}
	for i, test := range tests {
		require.Equal(t, test.expected, test.o.Architecture(), "test [%d] failed", i)
	}
}

func TestHome(t *testing.T) {
	type test struct {
		expected string
		o        os.OS
	}
	const (
		OtherHome = "/home/other"
	)
	tests := []test{
		{expected: os.MockUnixHomeDirectory, o: os.NewMock(os.WithPlatform(platform.Darwin))},
		{expected: os.MockUnixHomeDirectory, o: os.NewMock(os.WithPlatform(platform.Linux))},
		{expected: os.MockWindowsHomeDirectory, o: os.NewMock(os.WithPlatform(platform.Windows))},
		{expected: OtherHome, o: os.NewMock(os.WithHomeDirectory(OtherHome))},
	}
	for i, test := range tests {
		require.Equal(t, test.expected, test.o.Home(), "test [%d] failed", i)
	}
}

func TestWorkingDirectory(t *testing.T) {
	type test struct {
		expected string
		o        os.OS
	}
	const (
		OtherWorkingDirectory = "/home/other/wd"
	)
	tests := []test{
		{expected: os.MockUnixWorkingDirectory, o: os.NewMock(os.WithPlatform(platform.Darwin))},
		{expected: os.MockUnixWorkingDirectory, o: os.NewMock(os.WithPlatform(platform.Linux))},
		{expected: os.MockWindowsWorkingDirectory, o: os.NewMock(os.WithPlatform(platform.Windows))},
		{expected: OtherWorkingDirectory, o: os.NewMock(os.WithWorkingDirectory(OtherWorkingDirectory))},
		{expected: OtherWorkingDirectory, o: os.NewMock(os.WithWorkingDirectory(OtherWorkingDirectory))},
		{expected: OtherWorkingDirectory, o: os.NewMock(os.WithWorkingDirectory(OtherWorkingDirectory))},
		{expected: OtherWorkingDirectory, o: os.NewMock(os.WithWorkingDirectory(OtherWorkingDirectory))},
	}
	for i, test := range tests {
		workingDirectory, err := test.o.WorkingDirectory()
		require.Nil(t, err, "test [%d] o.WorkingDirectory() returned error")
		require.Equal(t, test.expected, workingDirectory, "test [%d] failed", i)
	}
}
