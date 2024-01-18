package message

import (
	"fmt"

	"github.com/jon4hz/gmotd/context"
	"github.com/jon4hz/gmotd/section/docker"
	"github.com/jon4hz/gmotd/section/hostname"
	"github.com/jon4hz/gmotd/section/smart"
	"github.com/jon4hz/gmotd/section/sysinfo"
	"github.com/jon4hz/gmotd/section/uptime"
	"github.com/jon4hz/gmotd/section/zpool"
)

type Section interface {
	fmt.Stringer

	Gather(*context.Context) error
	Print(*context.Context) string
}

var Message = []Section{
	hostname.Section{},
	uptime.Section{},
	sysinfo.Section{},
	zpool.Section{},
	docker.Section{},
	smart.Section{},
}
