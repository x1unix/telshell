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
	tprintln(rw, CLRF+decorations)
	tprintln(rw, banner)
	tprintln(rw, decorations+CLRF)
	return nil
}

// tprintln prints message with DOS line endings
func tprintln(rw io.ReadWriter, msg string) {
	fmt.Fprint(rw, msg+"\r\n")
}
