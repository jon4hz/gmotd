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
	Uptime   time.Duration
	Platform string
	Kernel   string
	CPU      string
	Load     *LoadInfo
}

type LoadInfo struct {
	Load1, Load5, Load15 float64
}
