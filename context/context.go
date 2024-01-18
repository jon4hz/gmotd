package context

import (
	"context"
	"time"

	"github.com/jon4hz/gmotd/config"
)

type Context struct {
	context.Context
	cancel context.CancelFunc

	Config   *config.Config
	Runtime  *Runtime
	Hostname *Hostname
	Uptime   *Uptime
	Sysinfo  *SysInfo
	Zpool    *Zpool
	Docker   *Docker
	Smart    *Smart
}

type Runtime struct {
	Width, Height int
}

func New() *Context {
	c, cancel := context.WithTimeout(context.Background(), time.Second*30)
	return &Context{
		Context: c,
		cancel:  cancel,
		Config:  &config.Config{},
		Runtime: &Runtime{},
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

type DiskInfo struct {
	Name        string
	Temperature uint64
	HasError    bool
}
