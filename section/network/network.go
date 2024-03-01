package network

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/jon4hz/gmotd/context"
	"github.com/shirou/gopsutil/v3/net"
)

type Section struct{}

func (Section) String() string { return "network" }

func (Section) Enabled(c *context.Context) bool {
	return !c.Config.Network.Disabled
}

/* func (Section) Default(ctx *context.Context) {
	ctx.Config.Uptime.Precision = 3
} */

func (Section) Gather(c *context.Context) error {
	stats, err := net.IOCountersWithContext(c, true)
	if err != nil {
		return err
	}

	for _, s := range stats {
		fmt.Println(s.Name)
		fmt.Println(humanize.Bytes(s.BytesRecv))
		fmt.Println(humanize.Bytes(s.BytesSent))

	}
	return nil
}

func (Section) Print(ctx *context.Context) string {
	c := ctx.Network
	if c == nil {
		return ""
	}
	return "network stat:"
}
