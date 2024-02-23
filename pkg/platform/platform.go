package platform

import (
	"context"

	"github.com/shirou/gopsutil/v3/host"
)

func defaultPrettyName(ctx context.Context) string {
	p, _, _, err := host.PlatformInformationWithContext(ctx)
	if err != nil {
		return ""
	}
	return p
}
