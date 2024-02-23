package plex

import (
	"fmt"

	"github.com/jon4hz/gmotd/context"
	"github.com/jon4hz/gmotd/pkg/plex"
	"github.com/jon4hz/gmotd/styles"
	"github.com/spf13/viper"
)

type Section struct{}

func (Section) String() string { return "plex" }

func (Section) Enabled(c *context.Context) bool {
	return !c.Config.Plex.Disabled
}

func (Section) Default(c *context.Context) {
	viper.SetDefault("plex.disabled", true)
	viper.SetDefault("plex.timeout", 5)
	viper.SetDefault("plex.tls_verify", true)
}

func (Section) Gather(c *context.Context) error {
	p := plex.New(c.Config.Plex.Server, c.Config.Plex.Token, c.Config.Plex.Timeout, c.Config.Plex.TLSVerify)

	sessions, err := p.CountSessions(c)
	if err != nil {
		return fmt.Errorf("failed to get plex sessions: %w", err)
	}

	c.Plex = &context.Plex{
		Sessions: sessions,
	}

	return nil
}

func (Section) Print(ctx *context.Context) string {
	c := ctx.Plex
	if c == nil {
		return ""
	}

	return fmt.Sprintf("plex sessions: %s", styles.Green.Render(fmt.Sprintf("%d", c.Sessions)))
}
