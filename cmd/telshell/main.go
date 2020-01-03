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
	shell := flag.String("shell", telshell.DefaultShell, "Shell to use")
	flag.Var(&shellFlags, "shellarg", "Shell args")
	flag.Parse()
	if err := start(*addr, *shell, shellFlags); err != nil {
		zap.S().Fatal(err)
	}
}

func start(addr, shell string, shellArgs []string) error {
	ctx, _ := app.GetApplicationContext()
	h := telshell.NewShellHandler(shell, shellArgs...)
	srv := telshell.NewServer(h)

	/*go func() {
		zap.S().Info("listening to ctx")
		for {
			select {
			case <-ctx.Done():
				zap.S().Info("ctx.Done()")
			default:
			}
		}
	}()*/

	return srv.Start(ctx, addr)
}
