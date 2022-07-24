package platform

import "github.com/shirou/gopsutil/v3/host"

func defaultPrettyName() string {
	p, _, _, err := host.PlatformInformation()
	if err != nil {

		return ""
	}
	return p
}
