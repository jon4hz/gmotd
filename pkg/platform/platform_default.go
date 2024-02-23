//go:build darwin || dragonfly || freebsd || netbsd || openbsd || solaris || windows
// +build darwin dragonfly freebsd netbsd openbsd solaris windows

package platform

import "context"

func PrettyName(ctx context.Context) string {
	return defaultPrettyName(ctx)
}
