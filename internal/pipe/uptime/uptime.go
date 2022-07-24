package uptime

import (
	"time"

	"github.com/hako/durafmt"
	"github.com/jon4hz/gmotd/internal/context"
	"github.com/shirou/gopsutil/host"
)

type Pipe struct{}

func (Pipe) String() string { return "uptime" }

func (Pipe) Message(c *context.Context) string {
	t, err := host.BootTime()
	if err != nil {
		return ""
	}
	uptime := time.Since(time.Unix(int64(t), 0)).Round(time.Second)
	return "Uptime: " + durafmt.Parse(uptime).LimitFirstN(2).String()
}
