package context

import "github.com/jon4hz/gmotd/internal/config"

type Context struct {
	Config  *config.Config
	Runtime *Runtime
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
