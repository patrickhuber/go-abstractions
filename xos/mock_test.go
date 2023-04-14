package xos_test

import (
	"testing"

	"github.com/patrickhuber/go-xplat/xos"
	"github.com/stretchr/testify/require"
)

func TestPlatform(t *testing.T) {
	type test struct {
		expected string
		o        xos.OS
	}
	tests := []test{
		{expected: xos.MockDarwinPlatform, o: xos.NewDarwinMock()},
		{expected: xos.MockLinuxPlatform, o: xos.NewLinuxMock()},
		{expected: xos.MockWindowsPlatform, o: xos.NewWindowsMock()},
		{expected: "plan9", o: xos.NewMock(xos.WithPlatform("plan9"))},
		{expected: "plan9", o: xos.NewDarwinMock(xos.WithPlatform("plan9"))},
		{expected: "plan9", o: xos.NewLinuxMock(xos.WithPlatform("plan9"))},
		{expected: "plan9", o: xos.NewWindowsMock(xos.WithPlatform("plan9"))},
	}
	for _, test := range tests {
		require.Equal(t, test.expected, test.o.Platform())
	}
}

func TestArchitecture(t *testing.T) {
	type test struct {
		expected string
		o        xos.OS
	}
	tests := []test{
		{expected: xos.MockAmd64Architecture, o: xos.NewDarwinMock()},
		{expected: xos.MockAmd64Architecture, o: xos.NewLinuxMock()},
		{expected: xos.MockAmd64Architecture, o: xos.NewWindowsMock()},
		{expected: xos.MockArm64Architecture, o: xos.NewMock(xos.WithArchitecture(xos.MockArm64Architecture))},
		{expected: xos.MockArm64Architecture, o: xos.NewWindowsMock(xos.WithArchitecture(xos.MockArm64Architecture))},
		{expected: xos.MockArm64Architecture, o: xos.NewLinuxMock(xos.WithArchitecture(xos.MockArm64Architecture))},
		{expected: xos.MockArm64Architecture, o: xos.NewDarwinMock(xos.WithArchitecture(xos.MockArm64Architecture))},
	}
	for i, test := range tests {
		require.Equal(t, test.expected, test.o.Architecture(), "test [%d] failed", i)
	}
}

func TestHome(t *testing.T) {
	type test struct {
		expected string
		o        xos.OS
	}
	const (
		OtherHome = "/home/other"
	)
	tests := []test{
		{expected: xos.MockDarwinHomeDirectory, o: xos.NewDarwinMock()},
		{expected: xos.MockLinuxHomeDirectory, o: xos.NewLinuxMock()},
		{expected: xos.MockWindowsHomeDirectory, o: xos.NewWindowsMock()},
		{expected: OtherHome, o: xos.NewMock(xos.WithHomeDirectory(OtherHome))},
		{expected: OtherHome, o: xos.NewWindowsMock(xos.WithHomeDirectory(OtherHome))},
		{expected: OtherHome, o: xos.NewLinuxMock(xos.WithHomeDirectory(OtherHome))},
		{expected: OtherHome, o: xos.NewDarwinMock(xos.WithHomeDirectory(OtherHome))},
	}
	for i, test := range tests {
		require.Equal(t, test.expected, test.o.Home(), "test [%d] failed", i)
	}
}

func TestWorkingDirectory(t *testing.T) {
	type test struct {
		expected string
		o        xos.OS
	}
	const (
		OtherWorkingDirectory = "/home/other/wd"
	)
	tests := []test{
		{expected: xos.MockDarwinWorkingDirectory, o: xos.NewDarwinMock()},
		{expected: xos.MockLinuxWorkingDirectory, o: xos.NewLinuxMock()},
		{expected: xos.MockWindowsWorkingDirectory, o: xos.NewWindowsMock()},
		{expected: OtherWorkingDirectory, o: xos.NewMock(xos.WithWorkingDirectory(OtherWorkingDirectory))},
		{expected: OtherWorkingDirectory, o: xos.NewWindowsMock(xos.WithWorkingDirectory(OtherWorkingDirectory))},
		{expected: OtherWorkingDirectory, o: xos.NewLinuxMock(xos.WithWorkingDirectory(OtherWorkingDirectory))},
		{expected: OtherWorkingDirectory, o: xos.NewDarwinMock(xos.WithWorkingDirectory(OtherWorkingDirectory))},
	}
	for i, test := range tests {
		workingDirectory, err := test.o.WorkingDirectory()
		require.Nil(t, err, "test [%d] o.WorkingDirectory() returned error")
		require.Equal(t, test.expected, workingDirectory, "test [%d] failed", i)
	}
}
