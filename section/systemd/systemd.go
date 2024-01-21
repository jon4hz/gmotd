package systemd

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/jon4hz/gmotd/context"
	"github.com/jon4hz/gmotd/styles"
)

type Section struct{}

func (Section) String() string { return "systemd" }

func (Section) Enabled(c *context.Context) bool {
	return !c.Config.Systemd.Disabled
}

func (Section) Gather(c *context.Context) error {
	conn, err := dbus.NewWithContext(c)
	if err != nil {
		return err
	}

	wanted := c.Config.Systemd.Units
	for i, unit := range wanted {
		if !strings.Contains(unit, ".") {
			wanted[i] = unit + ".service"
		}
	}

	units, err := conn.ListUnitsByNamesContext(c, wanted)
	if err != nil {
		return err
	}

	c.Systemd = &context.Systemd{
		Units: make([]context.SystemdUnit, len(units)),
	}
	for i, u := range units {
		c.Systemd.Units[i] = context.SystemdUnit{
			Name:        u.Name,
			ActiveState: u.ActiveState,
		}
	}

	return nil
}

var (
	keyStyle        = lipgloss.NewStyle()
	valueStyle      = lipgloss.NewStyle().Bold(true)
	leftColumnStyle = lipgloss.NewStyle().PaddingRight(4)
)

func (Section) Print(ctx *context.Context) string {
	c := ctx.Systemd
	if c == nil {
		return ""
	}

	units := make([]string, 0, len(c.Units))
	oddAdd := 0
	if cap(units)%2 != 0 {
		oddAdd = 1
	}

	units = append(units, renderColumn(c.Units[:len(c.Units)/2+oddAdd])...)
	units = append(units, renderColumn(c.Units[len(c.Units)/2+oddAdd:])...)

	leftColumn := lipgloss.JoinVertical(lipgloss.Top, units[:len(units)/2+oddAdd]...)
	bothColumns := lipgloss.JoinHorizontal(lipgloss.Left,
		leftColumnStyle.Render(leftColumn),
		lipgloss.JoinVertical(lipgloss.Top, units[len(units)/2+oddAdd:]...),
	)

	return lipgloss.JoinVertical(lipgloss.Top,
		"systemd status:",
		styles.Indent.Render(bothColumns),
	)
}

func renderColumn(u []context.SystemdUnit) []string {
	lnw := longestNameWidth(u)
	units := make([]string, len(u))
	for i, unit := range u {
		units[i] = renderUnit(unit, lnw)
	}
	return units
}

func renderUnit(unit context.SystemdUnit, lnw int) string {
	switch unit.ActiveState {
	case "failed":
		valueStyle.Foreground(styles.Red.GetForeground())
	case "inactive":
		valueStyle.Foreground(styles.Orange.GetForeground())
	default:
		valueStyle.Foreground(styles.Green.GetForeground())
	}

	return keyStyle.Width(lnw).Render(unit.Name+":") + "  " + valueStyle.Render(unit.ActiveState)
}
func longestNameWidth(units []context.SystemdUnit) int {
	l := 0
	for _, u := range units {
		if w := lipgloss.Width(u.Name) + 1; w > l {
			l = w
		}
	}
	return l
}
