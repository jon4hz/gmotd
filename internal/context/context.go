package context

import "github.com/jon4hz/gmotd/internal/config"

type Context struct {
	Config *config.Config
	Width  int
	Height int
}

func New() *Context {
	return &Context{
		Config: &config.Config{},
	}
}
