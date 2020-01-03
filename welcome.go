package telshell

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

type WelcomeHandler struct{}

func (h WelcomeHandler) Handle(_ context.Context, rw io.ReadWriter) error {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintln(rw, "ERROR: \t", err.Error())
		return err
	}

	banner := fmt.Sprintf("#  Wellcome to TelShell on %s (%s)  #", hostname, runtime.GOOS)
	decorations := strings.Repeat("#", len(banner))
	fmt.Fprintln(rw, "\n"+decorations)
	fmt.Fprintln(rw, banner)
	fmt.Fprintln(rw, decorations+"\n")
	return nil
}
