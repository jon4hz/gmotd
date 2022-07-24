package utils

import (
	"os"
	"path/filepath"
)

func HostEtc(combineWith ...string) string {
	return GetEnv("HOST_ETC", "/etc", combineWith...)
}

// GetEnv retrieves the environment variable key. If it does not exist it returns the default.
func GetEnv(key string, dfault string, combineWith ...string) string {
	value := os.Getenv(key)
	if value == "" {
		value = dfault
	}

	switch len(combineWith) {
	case 0:
		return value
	case 1:
		return filepath.Join(value, combineWith[0])
	default:
		all := make([]string, len(combineWith)+1)
		all[0] = value
		copy(all[1:], combineWith)
		return filepath.Join(all...)
	}
}

// PathExistsWithContents returns the filename exists and it is not empty
func PathExistsWithContents(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return info.Size() > 4 // at least 4 bytes
}
