package uptime

import (
	"fmt"
	"time"

	"github.com/hako/durafmt"
	"github.com/jon4hz/gmotd/context"
	"github.com/shirou/gopsutil/host"
)

type Pipe struct{}

func (Pipe) String() string { return "uptime" }

func (Pipe) Gather(c *context.Context) error {
	t, err := host.BootTime()
	if err != nil {
		return fmt.Errorf("failed to get uptime: %w", err)
	}
	c.Uptime = &context.Uptime{
		Uptime: time.Since(time.Unix(int64(t), 0)),
	}
	return nil
}

func (Pipe) Print(c *context.Context) string {
	if c.Uptime == nil {
		return ""
	}
	return "Uptime: " + durafmt.Parse(c.Uptime.Uptime).LimitFirstN(2).String()
}
