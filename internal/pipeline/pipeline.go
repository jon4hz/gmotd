package pipeline

import (
	"fmt"

	"github.com/jon4hz/gmotd/internal/context"
	"github.com/jon4hz/gmotd/internal/pipe/hostname"
	"github.com/jon4hz/gmotd/internal/pipe/sysinfo"
	"github.com/jon4hz/gmotd/internal/pipe/uptime"
)

type Pipe interface {
	fmt.Stringer
	Message(*context.Context) string
}

var Pipeline = []Pipe{
	hostname.Pipe{},
	uptime.Pipe{},
	sysinfo.Pipe{},
}
