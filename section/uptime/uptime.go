package uptime

import (
	"fmt"
	"time"

	"github.com/hako/durafmt"
	"github.com/jon4hz/gmotd/context"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/spf13/viper"
)

type Section struct{}

func (Section) String() string { return "uptime" }

func (Section) Enabled(c *context.Context) bool {
	return !c.Config.Uptime.Disabled
}

func (Section) Default(ctx *context.Context) {
	viper.SetDefault("uptime.precisioin", 3)
}

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
	return "uptime: " + durafmt.Parse(c.Uptime).LimitFirstN(ctx.Config.Uptime.Precision).String()
}
