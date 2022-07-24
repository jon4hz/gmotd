package sysinfo

import (
	"fmt"
	"strings"
	"time"

	"github.com/jon4hz/gmotd/context"

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
		Uptime: time.Since(time.Unix(int64(t), 0)),
	}
	return nil
}

func (Pipe) Print(c *context.Context) string {
	if c.Sysinfo == nil {
		return ""
	}
	var s strings.Builder
	s.WriteString(fmt.Sprintf("Uptime: %s\n", c.Uptime.Uptime))

	platform, family, version, err := host.PlatformInformation()
	if err != nil {
		return ""
	}
	s.WriteString(fmt.Sprintf("Platform: %s %s %s\n", platform, family, version))

	return s.String()
}
