package os_test

import (
	"testing"

	"github.com/patrickhuber/go-xplat/os"
	"github.com/stretchr/testify/require"
)

func TestPlatform(t *testing.T) {
	type test struct {
		expected string
		o        os.OS
	}
	tests := []test{
		{expected: os.MockDarwinPlatform, o: os.NewDarwinMock()},
		{expected: os.MockLinuxPlatform, o: os.NewLinuxMock()},
		{expected: os.MockWindowsPlatform, o: os.NewWindowsMock()},
		{expected: "plan9", o: os.NewMock(os.WithPlatform("plan9"))},
		{expected: "plan9", o: os.NewDarwinMock(os.WithPlatform("plan9"))},
		{expected: "plan9", o: os.NewLinuxMock(os.WithPlatform("plan9"))},
		{expected: "plan9", o: os.NewWindowsMock(os.WithPlatform("plan9"))},
	}
	for _, test := range tests {
		require.Equal(t, test.expected, test.o.Platform())
	}
}

func TestArchitecture(t *testing.T) {
	type test struct {
		expected string
		o        os.OS
	}
	tests := []test{
		{expected: os.MockAmd64Architecture, o: os.NewDarwinMock()},
		{expected: os.MockAmd64Architecture, o: os.NewLinuxMock()},
		{expected: os.MockAmd64Architecture, o: os.NewWindowsMock()},
		{expected: os.MockArm64Architecture, o: os.NewMock(os.WithArchitecture(os.MockArm64Architecture))},
		{expected: os.MockArm64Architecture, o: os.NewWindowsMock(os.WithArchitecture(os.MockArm64Architecture))},
		{expected: os.MockArm64Architecture, o: os.NewLinuxMock(os.WithArchitecture(os.MockArm64Architecture))},
		{expected: os.MockArm64Architecture, o: os.NewDarwinMock(os.WithArchitecture(os.MockArm64Architecture))},
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
		{expected: os.MockDarwinHomeDirectory, o: os.NewDarwinMock()},
		{expected: os.MockLinuxHomeDirectory, o: os.NewLinuxMock()},
		{expected: os.MockWindowsHomeDirectory, o: os.NewWindowsMock()},
		{expected: OtherHome, o: os.NewMock(os.WithHomeDirectory(OtherHome))},
		{expected: OtherHome, o: os.NewWindowsMock(os.WithHomeDirectory(OtherHome))},
		{expected: OtherHome, o: os.NewLinuxMock(os.WithHomeDirectory(OtherHome))},
		{expected: OtherHome, o: os.NewDarwinMock(os.WithHomeDirectory(OtherHome))},
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
		{expected: os.MockDarwinWorkingDirectory, o: os.NewDarwinMock()},
		{expected: os.MockLinuxWorkingDirectory, o: os.NewLinuxMock()},
		{expected: os.MockWindowsWorkingDirectory, o: os.NewWindowsMock()},
		{expected: OtherWorkingDirectory, o: os.NewMock(os.WithWorkingDirectory(OtherWorkingDirectory))},
		{expected: OtherWorkingDirectory, o: os.NewWindowsMock(os.WithWorkingDirectory(OtherWorkingDirectory))},
		{expected: OtherWorkingDirectory, o: os.NewLinuxMock(os.WithWorkingDirectory(OtherWorkingDirectory))},
		{expected: OtherWorkingDirectory, o: os.NewDarwinMock(os.WithWorkingDirectory(OtherWorkingDirectory))},
	}
	for i, test := range tests {
		workingDirectory, err := test.o.WorkingDirectory()
		require.Nil(t, err, "test [%d] o.WorkingDirectory() returned error")
		require.Equal(t, test.expected, workingDirectory, "test [%d] failed", i)
	}
}
