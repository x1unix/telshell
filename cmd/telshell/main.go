package main

import (
	"flag"
	"github.com/x1unix/telshell"
	"github.com/x1unix/telshell/internal/app"
	"go.uber.org/zap"
)

func init() {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(l)
}

func main() {
	shellFlags := telshell.ShellArgs[:]
	addr := flag.String("addr", ":1000", "Address to listen")
	withAuth := flag.Bool("auth", false, "Require authorization")
	shell := flag.String("shell", telshell.DefaultShell, "Define shell argument")
	flag.Var(&shellFlags, "s", "Shell args")

	flag.Parse()
	if err := start(*addr, *shell, *withAuth, shellFlags); err != nil {
		zap.S().Fatal(err)
	}
}

func start(addr, shell string, withAuth bool, shellArgs []string) error {
	ctx, _ := app.GetApplicationContext()
	var h telshell.Handler
	if withAuth {
		zap.L().Info("shell auth enabled")
		h = telshell.NewAuthShellHandler(shell, shellArgs...)
	} else {
		zap.S().Infof("shell auth disabled, shell is %q", shell)
		h = telshell.NewShellHandler(shell, shellArgs...)
	}
	srv := telshell.NewServer(telshell.WelcomeHandler{}, h)
	return srv.Start(ctx, addr)
}
