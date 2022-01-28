package sysinfo

import (
	"fmt"
	"strings"
	"time"

	"github.com/jon4hz/gmotd/internal/context"

	"github.com/shirou/gopsutil/v3/host"
)

type Pipe struct{}

func (Pipe) String() string { return "sysinfo" }

func (Pipe) Message(c *context.Context) string {
	var s strings.Builder

	t, err := host.BootTime()
	if err != nil {
		return ""
	}
	uptime := time.Since(time.Unix(int64(t), 0)).Round(time.Second)
	s.WriteString(fmt.Sprintf("Uptime: %s\n", uptime))

	platform, family, version, err := host.PlatformInformation()
	if err != nil {
		return ""
	}
	s.WriteString(fmt.Sprintf("Platform: %s %s %s\n", platform, family, version))

	return s.String()
}
