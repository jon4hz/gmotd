package diskspace

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/jon4hz/gmotd/config"
	"github.com/jon4hz/gmotd/context"
	"github.com/jon4hz/gmotd/styles"
	"github.com/muesli/termenv"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/spf13/viper"
)

type Section struct{}

func (Section) String() string { return "diskspace" }

func (Section) Enabled(c *context.Context) bool {
	return !c.Config.DiskSpace.Disabled
}

func (Section) Default(ctx *context.Context) {
	viper.SetDefault("diskspace.excluded_fs", []string{
		"tmpfs",
		"devtmpfs",
		"devfs",
		"iso9660",
		"overlay",
		"aufs",
		"squashfs",
		"zfs",
	})
}

func (Section) Gather(c *context.Context) error {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return err
	}

	disks := make([]context.Disk, 0)
	for _, p := range partitions {
		if fsIsExcluded(c.Config.DiskSpace, p.Fstype) {
			continue
		}
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			return err
		}
		d := context.Disk{
			Mountpoint: p.Mountpoint,
			Total:      usage.Total,
			Used:       usage.Used,
		}
		disks = append(disks, d)
	}

	sort.Slice(disks, func(i, j int) bool {
		return disks[i].Mountpoint < disks[j].Mountpoint
	})

	c.DiskSpace = &context.DiskSpace{
		Disks: disks,
	}

	return nil
}

func fsIsExcluded(cfg *config.DiskSpace, fs string) bool {
	for _, e := range cfg.ExcludedFS {
		if e == fs {
			return true
		}
	}
	return false
}

const progressWidth = 60

var statStyle = lipgloss.NewStyle().Align(lipgloss.Right)

func (Section) Print(ctx *context.Context) string {
	c := ctx.DiskSpace
	if c == nil {
		return ""
	}

	disks := make([]string, len(c.Disks)+1)
	for i, p := range c.Disks {
		disks[i+1] = styles.Indent.Render(renderDisk(p))
	}
	disks[0] = "disk usage:"

	return lipgloss.JoinVertical(lipgloss.Top, disks...)
}

func renderDisk(d context.Disk) string {
	used := float64(d.Used) / float64(d.Total)

	stat := statStyle.Width(progressWidth - lipgloss.Width(d.Mountpoint) + 1).
		Render(fmt.Sprintf(
			"%d%% used out of %s",
			uint64(used*100),
			humanize.IBytes(d.Total),
		))

	prog := progress.New(
		progress.WithGradient("#9ece6a", "#fd665f"),
		progress.WithWidth(progressWidth),
		progress.WithoutPercentage(),
		progress.WithColorProfile(termenv.TrueColor),
	)
	prog.Full = '='
	prog.Empty = '='

	return lipgloss.JoinVertical(lipgloss.Top,
		d.Mountpoint+stat,
		"["+prog.ViewAs(used)+"]",
	)
}
