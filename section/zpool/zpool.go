package zpool

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/jon4hz/gmotd/context"
	"github.com/jon4hz/gmotd/styles"
	"github.com/krystal/go-zfs"
)

type Section struct{}

func (Section) String() string { return "zpool" }

func (Section) Enabled(c *context.Context) bool {
	return !c.Config.Zpool.Disabled
}

/* func (Section) Default(ctx *context.Context) {
	ctx.Config.Zpool.Disabled = true
} */

func (Section) Gather(c *context.Context) error {
	z := zfs.New()

	pools, err := z.ListPools(c)
	if err != nil {
		return fmt.Errorf("failed to get zpool info: %w", err)
	}

	zpools := make([]context.ZpoolPool, len(pools))
	for i, p := range pools {
		h, _ := p.Health()
		c, _ := p.Capacity()
		f, _ := p.Free()
		a, _ := p.Allocated()
		zpools[i] = context.ZpoolPool{
			Name:      p.Name,
			Health:    h,
			Capacity:  c,
			Free:      f,
			Allocated: a,
		}
	}
	sort.Slice(zpools, func(i, j int) bool {
		return zpools[i].Name < zpools[j].Name
	})

	c.Zpool = &context.Zpool{
		Pools: zpools,
	}

	return nil
}

func (Section) Print(ctx *context.Context) string {
	c := ctx.Zpool
	if c == nil || len(c.Pools) == 0 {
		return ""
	}

	pools := make([]string, len(c.Pools)+1)
	for i, p := range c.Pools {
		pools[i+1] = styles.Indent.Render(renderZpool(p))
	}
	pools[0] = "zpool usage:"

	return lipgloss.JoinVertical(lipgloss.Top, pools...)
}

const progressWidth = 60

var statStyle = lipgloss.NewStyle().Align(lipgloss.Right)

func renderZpool(zpool context.ZpoolPool) string {
	name := zpool.Name
	switch zpool.Health {
	case zfs.HealthOnline:
		// do nothing
	case zfs.HealthDegraded:
		name = styles.Orange.Render(name + fmt.Sprintf(" (%s)", zpool.Health))
	default:
		name = styles.Red.Render(name + fmt.Sprintf(" (%s)", zpool.Health))
	}

	stat := statStyle.Width(progressWidth - lipgloss.Width(name) + 1).
		Render(fmt.Sprintf(
			"%d%% used out of %s",
			zpool.Capacity,
			humanize.IBytes(uint64(zpool.Allocated+zpool.Free)),
		))

	used := float64(zpool.Allocated) / float64(zpool.Allocated+zpool.Free)
	prog := progress.New(
		progress.WithGradient("#9ece6a", "#fd665f"),
		progress.WithWidth(progressWidth),
		progress.WithoutPercentage(),
	)
	prog.Full = '='
	prog.Empty = '='

	return lipgloss.JoinVertical(lipgloss.Top,
		name+stat,
		"["+prog.ViewAs(used)+"]",
	)
}
