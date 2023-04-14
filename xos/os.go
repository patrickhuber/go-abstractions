package os

import (
	"os"
	"runtime"
)

type OS interface {
	WorkingDirectory() (string, error)
	Executable() (string, error)
	Platform() string
	Architecture() string
	Home() string
}

type realOS struct {
}

func New() OS {
	return &realOS{}
}

func (o *realOS) WorkingDirectory() (string, error) {
	return os.Getwd()
}

func (o *realOS) Executable() (string, error) {
	return os.Executable()
}

func (o *realOS) Platform() string {
	return runtime.GOOS
}

func (o *realOS) Architecture() string {
	return runtime.GOARCH
}

func (o *realOS) Home() string {
	dir, _ := os.UserHomeDir()
	return dir
}
