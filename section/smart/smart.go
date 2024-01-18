package smart

import (
	"fmt"
	"log"

	"github.com/anatol/smart.go"
	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/gmotd/context"
	"github.com/jon4hz/gmotd/styles"
)

type Section struct{}

func (Section) String() string { return "smart" }

func (Section) Gather(c *context.Context) error {

	// tmp
	c.Config.Smart.Disks = []string{"/dev/sda", "/dev/sdb", "/dev/sdc", "/dev/sdd", "/dev/sde", "/dev/sdf", "/dev/sdg", "/dev/sdh"}

	disk := make([]context.DiskInfo, 0, len(c.Config.Smart.Disks))
	for _, d := range c.Config.Smart.Disks {
		dev, err := smart.Open(d)
		if err != nil {
			log.Printf("failed to open smart device: %s", err)
			continue
		}
		defer dev.Close()

		a, err := dev.ReadGenericAttributes()
		if err != nil {
			log.Printf("failed to read generic attributes: %s", err)
			continue
		}

		hasErr := false
		switch sm := dev.(type) {
		case *smart.SataDevice:
			stLog, err := sm.ReadSMARTSelfTestLog()
			if err != nil {
				log.Printf("failed to read smart self test log: %s", err)
				continue
			}
			for _, s := range stLog.Entry {
				if s.Status != 0 {
					hasErr = true
				}
			}
		}

		disk = append(disk, context.DiskInfo{
			Name:        d,
			Temperature: a.Temperature,
			HasError:    hasErr,
		})
	}

	c.Smart = &context.Smart{Disks: disk}

	return nil
}

const maxTemp = 40

var leftColumnStyle = lipgloss.NewStyle().PaddingRight(4)

func (Section) Print(ctx *context.Context) string {
	c := ctx.Smart
	if c == nil || len(c.Disks) == 0 {
		return ""
	}

	disks := make([]string, 0, len(c.Disks))
	oddAdd := 0
	if cap(disks)%2 != 0 {
		oddAdd = 1
	}

	for _, d := range c.Disks {
		disks = append(disks, renderDisk(d))
	}

	leftColumn := lipgloss.JoinVertical(lipgloss.Top, disks[:len(disks)/2+oddAdd]...)
	bothColumns := lipgloss.JoinHorizontal(lipgloss.Left,
		leftColumnStyle.Render(leftColumn),
		lipgloss.JoinVertical(lipgloss.Top,
			disks[len(disks)/2+oddAdd:]...,
		),
	)

	return lipgloss.JoinVertical(lipgloss.Top,
		"smart status:",
		styles.Indent.Render(bothColumns),
	)
}

func renderDisk(d context.DiskInfo) string {
	return fmt.Sprintf("%s:  %s | %s", d.Name, renderTemp(d.Temperature), renderError(d.HasError))
}

func renderTemp(t uint64) string {
	if t == 0 {
		return styles.Red.Render("--C")
	}
	if t > maxTemp {
		return styles.Red.Render(fmt.Sprintf("%dC", t))
	}
	return styles.Green.Render(fmt.Sprintf("%dC", t))
}

func renderError(e bool) string {
	if e {
		return styles.Red.Render("with errors")
	}
	return styles.Green.Render("without errors")
}
