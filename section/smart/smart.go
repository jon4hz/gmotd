package smart

import (
	"fmt"
	"log"

	"github.com/anatol/smart.go"
	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/gmotd/context"
	"github.com/jon4hz/gmotd/styles"
	"github.com/spf13/viper"
)

type Section struct{}

func (Section) String() string { return "smart" }

func (Section) Enabled(c *context.Context) bool {
	return !c.Config.Smart.Disabled
}

func (Section) Default(ctx *context.Context) {
	viper.SetDefault("smart.disks", []string{"/dev/sda"})
}

func (Section) Gather(c *context.Context) error {
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
		var diskType context.DiskType
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

			iden, err := sm.Identify()
			if err != nil {
				log.Printf("failed to read identify: %s", err)
				continue
			}

			if iden.RotationRate != 0 {
				diskType = context.DiskTypeHDD
			} else {
				diskType = context.DiskTypeSSD
			}

		case *smart.NVMeDevice:
			diskType = context.DiskTypeNVMe
		}

		disk = append(disk, context.DiskInfo{
			Name:        d,
			Temperature: a.Temperature,
			HasError:    hasErr,
			Type:        diskType,
		})
	}

	c.Smart = &context.Smart{Disks: disk}

	return nil
}

const (
	maxTempHDD = 45
	maxTempSSD = 70
)

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

	lnw := longestNameWidth(c.Disks)

	for _, d := range c.Disks {
		disks = append(disks, renderDisk(d, lnw))
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

func renderDisk(d context.DiskInfo, lnw int) string {
	return fmt.Sprintf("%s%s | %s", renderName(d.Name, lnw), renderTemp(d.Temperature, d.Type), renderError(d.HasError))
}

func longestNameWidth(disks []context.DiskInfo) int {
	max := 0
	for _, d := range disks {
		if len(d.Name) > max {
			max = len(d.Name)
		}
	}
	return max
}

func renderName(name string, lnw int) string {
	return lipgloss.NewStyle().Width(lnw + 2).Render(name + ":")
}

func renderTemp(t uint64, dtype context.DiskType) string {
	if t == 0 {
		return styles.Red.Render("--C")
	}

	var maxTmp uint64
	switch dtype {
	case context.DiskTypeHDD:
		maxTmp = maxTempHDD
	case context.DiskTypeSSD, context.DiskTypeNVMe:
		maxTmp = maxTempSSD
	}

	if t > maxTmp {
		return styles.Red.Render(fmt.Sprintf("%dC", t))
	}
	return styles.Green.Render(fmt.Sprintf("%dC", t))
}

func renderError(e bool) string {
	if e {
		return styles.Red.Render("with errors!!!")
	}
	return styles.Green.Render("without errors")
}
