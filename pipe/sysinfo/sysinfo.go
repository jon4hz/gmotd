package sysinfo

import (
	"fmt"
	"strings"
	"time"

	"github.com/jon4hz/gmotd/context"
	"github.com/jon4hz/gmotd/internal/platform"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
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

	l, err := load.AvgWithContext(c)
	if err != nil {
		return fmt.Errorf("failed to get load info: %w", err)
	}

	c.Sysinfo = &context.SysInfo{
		Uptime:   time.Since(time.Unix(int64(s.BootTime), 0)),
		Platform: platform.PrettyName(c),
		Kernel:   s.KernelVersion,
		CPU:      cpuName,
		Load: &context.LoadInfo{
			Load1:  l.Load1,
			Load5:  l.Load5,
			Load15: l.Load15,
		},
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
	s.WriteString(fmt.Sprintf("Load: %.2f (1m), %.2f (5m), %.2f (15m)\n", c.Sysinfo.Load.Load1, c.Sysinfo.Load.Load5, c.Sysinfo.Load.Load15))
	return s.String()
}
