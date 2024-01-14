package sysinfo

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	humanize "github.com/dustin/go-humanize"
	"github.com/hako/durafmt"
	"github.com/jon4hz/gmotd/context"
	"github.com/jon4hz/gmotd/internal/platform"
	"github.com/jon4hz/gmotd/styles"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

type Section struct{}

func (Section) String() string { return "sysinfo" }

func (Section) Gather(c *context.Context) error {
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

	ps, err := process.Processes()
	if err != nil {
		return fmt.Errorf("failed to get process info: %w", err)
	}
	rootProcs := 0
	for _, p := range ps {
		uids, err := p.Uids()
		if err != nil {
			continue
		}
		for _, u := range uids {
			if u == 0 {
				rootProcs++
				break
			}
		}
	}

	memory, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("failed to get memory info: %w", err)
	}

	c.Sysinfo = &context.SysInfo{
		Uptime:   time.Since(time.Unix(int64(s.BootTime), 0)),
		Platform: platform.PrettyName(c),
		Kernel:   s.KernelVersion,
		CPU:      cpuName,
		CPUCount: len(cpu),
		Load: &context.LoadInfo{
			Load1:  l.Load1,
			Load5:  l.Load5,
			Load15: l.Load15,
		},
		RootProcs:   rootProcs,
		UserProcs:   len(ps) - rootProcs,
		MemoryTotal: memory.Total,
		MemoryUsed:  memory.Used,
		MemoryFree:  memory.Free,
	}
	return nil
}

var (
	style = func(s string) string {
		return lipgloss.PlaceHorizontal(14, lipgloss.Top, styles.Indent.Render(s), lipgloss.WithWhitespaceChars("."))
	}

	greenF = func(x float64) string {
		return styles.Green.Render(fmt.Sprintf("%.2f", x))
	}
	green = func(x int) string {
		return styles.Green.Render(fmt.Sprintf("%d", x))
	}
)

func (Section) Print(ctx *context.Context) string {
	c := ctx.Sysinfo
	if c == nil {
		return ""
	}

	return lipgloss.JoinVertical(lipgloss.Top,
		"system info:",
		style("Distro")+": "+c.Platform,
		style("Kernel")+": "+c.Kernel,
		"",
		style("Uptime")+": "+durafmt.Parse(c.Uptime).LimitFirstN(3).String(),
		style("Load")+":"+fmt.Sprintf(" %s (1m), %s (5m), %s (15m)", greenF(c.Load.Load1), greenF(c.Load.Load5), greenF(c.Load.Load15)),
		style("Processes")+":"+fmt.Sprintf(" %s (root), %s (user), %s (total)", green(c.RootProcs), green(c.UserProcs), green(c.RootProcs+c.UserProcs)),
		"",
		style("CPU")+": "+c.CPU+fmt.Sprintf(" (%s cores)", green(c.CPUCount)),

		style("Memory")+": "+fmt.Sprintf("%s used, %s free, %s total",
			styles.ColorizeWithMax(c.MemoryUsed, c.MemoryTotal).Render(humanize.IBytes(c.MemoryUsed)),
			styles.ColorizeWithMin(c.MemoryFree, c.MemoryTotal).Render(humanize.IBytes(c.MemoryFree)),
			styles.Green.Render(humanize.IBytes(c.MemoryTotal))),
	)
}
