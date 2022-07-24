//go:build linux
// +build linux

package platform

import (
	"github.com/jon4hz/gmotd/utils"
	"gopkg.in/ini.v1"
)

func PrettyName() string {
	if utils.PathExistsWithContents(utils.HostEtc("os-release")) {
		cfg, err := ini.Load(utils.HostEtc("os-release"))
		if err != nil {
			return defaultPrettyName()
		}
		return cfg.Section("").Key("PRETTY_NAME").String()
	}

	return defaultPrettyName()
}
