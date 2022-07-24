package defaults

import (
	"fmt"

	"github.com/jon4hz/gmotd/context"
	"github.com/jon4hz/gmotd/pipe/hostname"
)

type Defaulter interface {
	fmt.Stringer

	Default(ctx *context.Context)
}

var Defaulters = []Defaulter{
	hostname.Pipe{},
}
