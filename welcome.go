package telshell

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
)

type WelcomeHandler struct{}

func (h WelcomeHandler) Handle(_ context.Context, rw io.ReadWriter) error {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintln(rw, "ERROR: \t", err.Error())
		return err
	}

	fmt.Fprintf(rw, "Wellcome to TelShell on %s (%s)\r\n", hostname, runtime.GOOS)
	return nil
}
