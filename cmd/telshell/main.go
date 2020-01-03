package main

import (
	"flag"
	"fmt"
	"github.com/x1unix/telshell"
	"github.com/x1unix/telshell/internal/app"
	"go.uber.org/zap"
	"os"
)

func init() {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(l)
}

type startupParams struct {
	addr       string
	shell      string
	shellArgs  app.FlagsArray
	withAuth   bool
	bufferSize int
}

func main() {
	params := startupParams{
		shellArgs: telshell.ShellArgs[:],
	}

	flag.StringVar(&params.addr, "addr", ":1000", "Address to listen")
	flag.BoolVar(&params.withAuth, "auth", false, "Require authorization")
	flag.StringVar(&params.shell, "shell", telshell.DefaultShell, "Define shell argument")
	flag.IntVar(&params.bufferSize, "buffer", 64, "Buffer size")
	flag.Var(&params.shellArgs, "s", "Shell args")

	flag.Usage = func() {
		fmt.Println("TelShell - Simple telnet shell server")
		fmt.Printf("\nUsage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if err := start(params); err != nil {
		zap.S().Fatal(err)
	}
}

func start(p startupParams) error {
	ctx, _ := app.GetApplicationContext()
	var h telshell.Handler
	if p.withAuth {
		zap.S().Infof("shell auth enabled, shell is %q", p.shell)
		h = telshell.NewAuthShellHandler(p.bufferSize, p.shell, p.shellArgs...)
	} else {
		zap.S().Infof("shell auth disabled, shell is %q", p.shell)
		h = telshell.NewShellHandler(p.bufferSize, p.shell, p.shellArgs...)
	}
	srv := telshell.NewServer(telshell.WelcomeHandler{}, h)
	return srv.Start(ctx, p.addr)
}
