// +build !windows,!js,!nacl

package telshell

import "github.com/x1unix/telshell/internal/app"

var (
	DefaultShell = "/bin/sh"
	ShellArgs = app.FlagsArray{"-i"}
)
