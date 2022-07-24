package pipeline

import (
	"fmt"

	"github.com/jon4hz/gmotd/context"
	"github.com/jon4hz/gmotd/pipe/hostname"
	"github.com/jon4hz/gmotd/pipe/sysinfo"
	"github.com/jon4hz/gmotd/pipe/uptime"
)

type Pipe interface {
	fmt.Stringer

	Gather(*context.Context) error
	Print(*context.Context) string
}

var Pipeline = []Pipe{
	hostname.Pipe{},
	uptime.Pipe{},
	sysinfo.Pipe{},
}
