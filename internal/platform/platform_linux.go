//go:build linux
// +build linux

package platform

import (
	"context"

	"github.com/jon4hz/gmotd/utils"
	"gopkg.in/ini.v1"
)

func PrettyName(ctx context.Context) string {
	if utils.PathExistsWithContents(utils.HostEtc("os-release")) {
		cfg, err := ini.Load(utils.HostEtc("os-release"))
		if err != nil {
			return defaultPrettyName(ctx)
		}
		return cfg.Section("").Key("PRETTY_NAME").String()
	}

	return defaultPrettyName(ctx)
}
