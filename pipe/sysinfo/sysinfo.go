package sysinfo

import (
	"fmt"
	"strings"
	"time"

	"github.com/jon4hz/gmotd/context"
	"github.com/jon4hz/gmotd/internal/platform"

	"github.com/shirou/gopsutil/v3/host"
)

type Pipe struct{}

func (Pipe) String() string { return "sysinfo" }

func (Pipe) Gather(c *context.Context) error {
	t, err := host.BootTime()
	if err != nil {
		return fmt.Errorf("failed to get uptime: %w", err)
	}

	c.Sysinfo = &context.Sysinfo{
		Uptime:   time.Since(time.Unix(int64(t), 0)),
		Platform: platform.PrettyName(),
	}
	return nil
}

func (Pipe) Print(c *context.Context) string {
	if c.Sysinfo == nil {
		return ""
	}
	var s strings.Builder
	s.WriteString(fmt.Sprintf("Uptime: %s\n", c.Sysinfo.Uptime))
	s.WriteString(fmt.Sprintf("Platform: %s\n", c.Sysinfo.Platform))
	return s.String()
}
