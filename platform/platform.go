// platform provides a runtime representation of the current platform
// see https://go.dev/doc/install/source#environment
package platform

import "runtime"

// Platform represents the specific platform information. Strings in platform are the values derived from runtime.GOOS
type Platform string

const (
	AIX       Platform = "aix"
	Android   Platform = "android"
	Darwin    Platform = "darwin"
	Dragonfly Platform = "dragonfly"
	FreeBSD   Platform = "freebsd"
	Hurd      Platform = "hurd"
	Illumos   Platform = "illumos"
	IOS       Platform = "ios"
	JS        Platform = "js"
	Linux     Platform = "linux"
	NACL      Platform = "nacl"
	NetBSD    Platform = "netbsd"
	OpenBSD   Platform = "openbsd"
	Plan9     Platform = "plan9"
	Solaris   Platform = "solaris"
	Wasip1    Platform = "wasip1"
	Windows   Platform = "windows"
	ZOS       Platform = "zos"
)

// IsUnix returns true if the platform is a unix platform
func (p Platform) IsUnix() bool {
	switch p {
	case AIX:
		return true
	case Android:
		return true
	case Darwin:
		return true
	case Dragonfly:
		return true
	case FreeBSD:
		return true
	case Hurd:
		return true
	case Illumos:
		return true
	case IOS:
		return true
	case Linux:
		return true
	case NetBSD:
		return true
	case OpenBSD:
		return true
	case Solaris:
		return true
	}
	return false
}

// IsWindows returns true if the platform is windows
func (p Platform) IsWindows() bool {
	return p == Windows
}

// String returns the string representation of the Platform
func (p Platform) String() string {
	return string(p)
}

// Default returns the default platform for the current runtime
func Default() Platform {
	return Platform(runtime.GOOS)
}
