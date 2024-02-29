package context

import (
	"context"
	"time"

	"github.com/jon4hz/gmotd/config"
)

type Context struct {
	context.Context
	cancel context.CancelFunc

	Config    *config.Config
	Hostname  *Hostname
	Uptime    *Uptime
	Sysinfo   *SysInfo
	Zpool     *Zpool
	DiskSpace *DiskSpace
	Docker    *Docker
	Smart     *Smart
	Systemd   *Systemd
	Plex      *Plex
}

/* type Runtime struct {
	Width, Height int
} */

func New() *Context {
	c, cancel := context.WithTimeout(context.Background(), time.Second*30)
	return &Context{
		Context: c,
		cancel:  cancel,
		Config: &config.Config{
			Hostname:  new(config.Hostname),
			Uptime:    new(config.Uptime),
			SysInfo:   new(config.SysInfo),
			Zpool:     new(config.Zpool),
			DiskSpace: new(config.DiskSpace),
			Docker:    new(config.Docker),
			Smart:     new(config.Smart),
			Systemd:   new(config.Systemd),
			Plex:      new(config.Plex),
		},
	}
}

func (c *Context) Cancel() {
	c.cancel()
}

type Hostname struct {
	Hostname string
}
type Uptime struct {
	Uptime time.Duration
}

type SysInfo struct {
	Uptime      time.Duration
	Platform    string
	Kernel      string
	CPU         string
	CPUCount    int
	Load        *LoadInfo
	RootProcs   int
	UserProcs   int
	MemoryTotal uint64
	MemoryUsed  uint64
	MemoryFree  uint64
}

type LoadInfo struct {
	Load1, Load5, Load15 float64
}

type Zpool struct {
	Pools []ZpoolPool
}

type ZpoolPool struct {
	Name      string
	Health    string
	Capacity  uint64
	Free      uint64
	Allocated uint64
}

type DiskSpace struct {
	Disks []Disk
}

type Disk struct {
	Mountpoint string
	Total      uint64
	Used       uint64
}

type Docker struct {
	Containers []DockerContainer
}

type DockerContainer struct {
	Name    string
	State   string
	Healthy bool
}

type Smart struct {
	Disks []DiskInfo
}

type DiskType string

const (
	DiskTypeHDD  DiskType = "hdd"
	DiskTypeSSD  DiskType = "ssd"
	DiskTypeNVMe DiskType = "nvme"
)

type DiskInfo struct {
	Name        string
	Type        DiskType
	Temperature uint64
	HasError    bool
}

type Systemd struct {
	Units []SystemdUnit
}

type SystemdUnit struct {
	Name        string
	ActiveState string
}

type Plex struct {
	Sessions int
}
