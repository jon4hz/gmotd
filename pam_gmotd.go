package main

import (
	"fmt"
	"io"
	"log/syslog"
	"os"
	"strings"
)

type SessionResult int

const (
	SessionOK SessionResult = iota
	SessionError
)

func pamLog(format string, args ...interface{}) {
	l, err := syslog.New(syslog.LOG_AUTH|syslog.LOG_WARNING, "pam-gmotd")
	if err != nil {
		return
	}
	l.Warning(fmt.Sprintf(format, args...)) // nolint:errcheck
}

func session(w io.Writer, uid int, configFile string) SessionResult {
	pamLog("session!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")

	origEUID := os.Geteuid()
	if os.Getuid() != origEUID || origEUID == 0 {
		// Note: this only sets the euid and doesn't do anything with the egid.
		// That should be fine for most cases, but it's worth calling out.
		if !seteuid(uid) {
			pamLog("error dropping privs from %d to %d", origEUID, uid)
			return SessionError
		}
		defer func() {
			if !seteuid(origEUID) {
				pamLog("error resetting uid to %d", origEUID)
			}
		}()
	}

	fmt.Fprintf(w, "config: %s\n", configFile)
	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

	return SessionOK
}

func pamSession(w io.Writer, uid int, username string, argv []string) SessionResult {
	configFile := "/etc/gmotd/gmotd.yml"

	for _, arg := range argv {
		opt := strings.Split(arg, "=")
		switch opt[0] {
		case "config":
			configFile = opt[1]
			pamLog("config set to %s", configFile)
		default:
			pamLog("unkown option: %s\n", opt[0])
		}
	}
	return session(w, uid, configFile)
}

func main() {}
