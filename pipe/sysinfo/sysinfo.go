package sysinfo

import (
	"fmt"
	"strings"
	"time"

	"github.com/jon4hz/gmotd/context"
	"github.com/jon4hz/gmotd/internal/platform"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
)

type Pipe struct{}

func (Pipe) String() string { return "sysinfo" }

func (Pipe) Gather(c *context.Context) error {
	s, err := host.InfoWithContext(c)
	if err != nil {
		return fmt.Errorf("failed to get host info: %w", err)
	}

	cpu, err := cpu.InfoWithContext(c)
	if err != nil {
		return fmt.Errorf("failed to get cpu info: %w", err)
	}

	var cpuName string
	if len(cpu) > 0 {
		cpuName = cpu[0].ModelName
	}

	c.Sysinfo = &context.Sysinfo{
		Uptime:   time.Since(time.Unix(int64(s.BootTime), 0)),
		Platform: platform.PrettyName(c),
		Kernel:   s.KernelVersion,
		CPU:      cpuName,
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
	s.WriteString(fmt.Sprintf("Kernel: %s\n", c.Sysinfo.Kernel))
	s.WriteString(fmt.Sprintf("CPU: %s\n", c.Sysinfo.CPU))
	return s.String()
}
