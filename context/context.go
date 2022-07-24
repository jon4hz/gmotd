package context

import (
	"time"

	"github.com/jon4hz/gmotd/config"
)

type Context struct {
	Config   *config.Config
	Runtime  *Runtime
	Hostname *Hostname
	Uptime   *Uptime
	Sysinfo  *Sysinfo
}

type Runtime struct {
	Width, Height int
}

func New() *Context {
	return &Context{
		Config:  &config.Config{},
		Runtime: &Runtime{},
	}
}

type Hostname struct {
	Hostname string
}
type Uptime struct {
	Uptime time.Duration
}

type Sysinfo struct {
	Uptime   time.Duration
	Platform string
	Kernel   string
	CPU      string
}
