package uptime

import (
	"fmt"
	"time"

	"github.com/hako/durafmt"
	"github.com/jon4hz/gmotd/context"
	"github.com/shirou/gopsutil/v3/host"
)

type Section struct{}

func (Section) String() string { return "uptime" }

func (Section) Gather(c *context.Context) error {
	t, err := host.BootTimeWithContext(c)
	if err != nil {
		return fmt.Errorf("failed to get uptime: %w", err)
	}
	c.Uptime = &context.Uptime{
		Uptime: time.Since(time.Unix(int64(t), 0)),
	}
	return nil
}

func (Section) Print(ctx *context.Context) string {
	c := ctx.Uptime
	if c == nil {
		return ""
	}
	return "uptime: " + durafmt.Parse(c.Uptime).LimitFirstN(3).String()
}