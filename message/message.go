package message

import (
	"fmt"

	"github.com/jon4hz/gmotd/context"
	"github.com/jon4hz/gmotd/section/diskspace"
	"github.com/jon4hz/gmotd/section/docker"
	"github.com/jon4hz/gmotd/section/hostname"
	"github.com/jon4hz/gmotd/section/plex"
	"github.com/jon4hz/gmotd/section/smart"
	"github.com/jon4hz/gmotd/section/sysinfo"
	"github.com/jon4hz/gmotd/section/systemd"
	"github.com/jon4hz/gmotd/section/uptime"
	"github.com/jon4hz/gmotd/section/zpool"
)

type Defaulter interface {
	fmt.Stringer

	Default(ctx *context.Context)
}

type Section interface {
	fmt.Stringer

	Enabled(*context.Context) bool
	Gather(*context.Context) error
	Print(*context.Context) string
}

var Message = []Section{
	hostname.Section{},
	uptime.Section{},
	sysinfo.Section{},
	zpool.Section{},
	diskspace.Section{},
	smart.Section{},
	docker.Section{},
	systemd.Section{},
	plex.Section{},
}
