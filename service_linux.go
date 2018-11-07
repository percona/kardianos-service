// Copyright 2015 Daniel Theophanes.
// Use of this source code is governed by a zlib-style
// license that can be found in the LICENSE file.

package service

import (
	"os"
	"strings"
)

type linuxSystemService struct {
	name        string
	detect      func() bool
	interactive func() bool
	new         func(i Interface, c *Config) (Service, error)
}

func (sc linuxSystemService) String() string {
	return sc.name
}
func (sc linuxSystemService) Detect() bool {
	return sc.detect()
}
func (sc linuxSystemService) Interactive() bool {
	return sc.interactive()
}
func (sc linuxSystemService) New(i Interface, c *Config) (Service, error) {
	return sc.new(i, c)
}

func init() {
	ChooseSystem(
		linuxSystemService{
			name:   "linux-systemd",
			detect: isSystemd,
			interactive: func() bool {
				is, _ := isInteractive()
				return is
			},
			new: newSystemdService,
		},
		linuxSystemService{
			name:   "linux-upstart",
			detect: isUpstart,
			interactive: func() bool {
				is, _ := isInteractive()
				return is
			},
			new: newUpstartService,
		},
		linuxSystemService{
			name:   "linux-supervisord",
			detect: isSupervisord,
			interactive: func() bool {
				is, _ := isInteractive()
				return is
			},
			new: newSupervisordService,
		},
		linuxSystemService{
			name:   "unix-systemv",
			detect: func() bool { return true },
			interactive: func() bool {
				is, _ := isInteractive()
				return is
			},
			new: newSystemVService,
		},
	)
}

func isInteractive() (bool, error) {
	// TODO: This is not true for user services.
	return os.Getppid() != 1, nil
}

var tf = map[string]interface{}{
	"cmd": func(s string) string {
		// Put command in single quotes, otherwise special characters like dollar ($) sign will be interpreted.
		return `'` + strings.Replace(s, `'`, `'"'"'`, -1) + `'`
	},
	"cmdSystemD": func(s string) string {
		s = strings.Replace(s, `%`, `%%`, -1)
		s = `"` + strings.Replace(s, `"`, `\"`, -1) + `"`
		return s
	},
	"cmdEscape": func(s string) string {
		return strings.Replace(s, " ", `\x20`, -1)
	},
	"envKey": func(env string) string {
		return strings.Split(env, "=")[0]
	},
	"envValue": func(env string) string {
		return strings.Join(strings.Split(env, "=")[1:], "=")
	},
}
