//go:build darwin || dragonfly || freebsd || netbsd || openbsd || solaris || windows
// +build darwin dragonfly freebsd netbsd openbsd solaris windows

package platform

func PrettyName() string {
	return defaultPrettyName()
}
