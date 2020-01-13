package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/x1unix/telshell"
	"github.com/x1unix/telshell/internal/app"
	"go.uber.org/zap"
)

const version = "1.1.0"

type startupParams struct {
	debug              bool
	withAuth           bool
	replaceLineEndings bool
	bufferSize         int
	addr               string
	shell              string
	shellArgs          app.FlagsArray
}

func (s startupParams) ioParams() telshell.IOParams {
	return telshell.IOParams{
		BufferSize:         uint(s.bufferSize),
		ReplaceLineEndings: s.replaceLineEndings,
	}
}

func (s startupParams) getLogger() (*zap.Logger, error) {
	if s.debug {
		return zap.NewDevelopment()
	}

	// Disable stacktrace and caller
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.CallerKey = ""
	cfg.EncoderConfig.NameKey = ""
	cfg.EncoderConfig.StacktraceKey = ""
	return cfg.Build()
}

func main() {
	params := startupParams{
		shellArgs: telshell.ShellArgs[:],
	}

	flag.StringVar(&params.addr, "addr", ":1000", "Address to listen")
	flag.BoolVar(&params.withAuth, "auth", false, "Require authorization")
	flag.BoolVar(&params.debug, "debug", false, "Enable debug mode")
	flag.StringVar(&params.shell, "shell", telshell.DefaultShell, "Define shell argument")
	flag.IntVar(&params.bufferSize, "buffer", 64, "Buffer size")
	flag.Var(&params.shellArgs, "s", "Define shell argument")
	flag.BoolVar(&params.replaceLineEndings, "replaceLineEndings", true,
		"Replace UNIX (\\n) with DOS (\\r\\n) line endings")

	flag.Usage = func() {
		fmt.Println("TelShell, version", version)
		fmt.Println("Simple telnet shell server")
		fmt.Printf("\nUsage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	// Initialize logger
	l, err := params.getLogger()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	zap.ReplaceGlobals(l)
	defer l.Sync()
	if err := start(params); err != nil {
		zap.S().Fatal(err)
	}
}

func start(p startupParams) error {
	ctx, _ := app.GetApplicationContext()
	var h telshell.Handler
	if p.withAuth {
		zap.S().Infof("shell auth enabled, shell is %q", p.shell)
		h = telshell.NewAuthShellHandler(p.ioParams(), p.shell, p.shellArgs...)
	} else {
		zap.S().Infof("shell auth disabled, shell is %q", p.shell)
		h = telshell.NewShellHandler(p.ioParams(), p.shell, p.shellArgs...)
	}
	srv := telshell.NewServer(telshell.WelcomeHandler{}, h)
	return srv.Start(ctx, p.addr)
}
