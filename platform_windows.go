package telshell

import "github.com/x1unix/telshell/internal/app"

var (
	DefaultShell = "cmd.exe"
	ShellArgs    = app.FlagsArray{}
)

func runShellAs(ctx context.Context, username, shell string, shellArgs ...string) *exec.Cmd {
	args := append([]string{shell}, shellArgs...)
	return exec.CommandContext(ctx, "runas", "/user:"+username, strings.Join(args, " "))
}
