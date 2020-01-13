// +build !windows,!js,!nacl

package telshell

import (
	"context"
	"github.com/x1unix/telshell/internal/app"
	"os"
	"os/exec"
)

var (
	DefaultShell = "/bin/sh"
	ShellArgs    = app.FlagsArray{}
)

func runShellAs(ctx context.Context, username, shell string, shellArgs ...string) *exec.Cmd {
	args := append([]string{"-k", "-Su", username, shell}, shellArgs...)
	cmd := exec.CommandContext(ctx, "sudo", args...)
	cmd.Env = os.Environ()
	return cmd
}
