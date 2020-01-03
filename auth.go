package telshell

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const MaxLoginChars = 32

var nameRegex = regexp.MustCompile(`(?m)^[A-Za-z0-9\.\-\_]+$`)

// AuthShellHandler is shell handler that requires username and password
type AuthShellHandler struct {
	buffSize  int
	shellPath string
	shellArgs []string
	log       *zap.SugaredLogger
}

// NewAuthShellHandler creates new authorized shell handler
func NewAuthShellHandler(buffSize int, shell string, args ...string) AuthShellHandler {
	return AuthShellHandler{
		shellPath: shell,
		shellArgs: args,
		buffSize:  buffSize,
		log:       zap.S().Named("auth_shell"),
	}
}

// Handle implements telshell.Handler
func (h AuthShellHandler) Handle(ctx context.Context, rw io.ReadWriter) error {
	fmt.Fprint(rw, "Username:")

	// Read output and match string
	buff := make([]byte, MaxLoginChars*4)
	_, _ = rw.Read(buff)

	// Sanitize input
	buff = TrimCLRF(buff)
	if len(buff) == 0 || !nameRegex.Match(buff) {
		// Ignore empty prompt response or invalid username format
		fmt.Fprintln(rw, "Access denied")
		return nil
	}

	username := string(buff)
	return h.startUserShell(ctx, username, rw)
}

func (h AuthShellHandler) startUserShell(ctx context.Context, user string, rw io.ReadWriter) error {
	wrapCtx, cancelFn := context.WithCancel(ctx)
	wrapper := NewTerminalWrapper(h.log, rw, h.buffSize)
	cmd := runShellAs(ctx, user, h.shellPath, h.shellArgs...)
	cmd.Env = os.Environ()
	if err := wrapper.Listen(wrapCtx, cmd); err != nil {
		return err
	}

	h.log.Debugw("login shell start",
		"command", cmd.Path,
		"args", cmd.Args,
	)
	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "failed to start shell instance")
	}

	defer cancelFn()
	return cmd.Wait()
}
