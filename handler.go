package telshell

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Handler is TCP request handler
type Handler interface {
	Handle(ctx context.Context, rw io.ReadWriter) error
}

// ShellHandler provides shell access via Telnet
type ShellHandler struct {
	buffSize  int
	shellPath string
	shellArgs []string
	log       *zap.SugaredLogger
}

// NewShellHandler creates a new shell handler
func NewShellHandler(buffSize int, shell string, args ...string) ShellHandler {
	return ShellHandler{
		shellPath: shell,
		shellArgs: args,
		buffSize:  buffSize,
		log:       zap.S().Named("shell"),
	}
}

// Handle implements telshell.Handler
func (s ShellHandler) Handle(ctx context.Context, rw io.ReadWriter) error {
	fmt.Fprintf(rw, "Current shell is %q\n", s.shellPath)

	wrapCtx, cancelFn := context.WithCancel(ctx)
	wrapper := NewTerminalWrapper(s.log, rw, s.buffSize)
	cmd := exec.CommandContext(ctx, s.shellPath, s.shellArgs...)
	cmd.Env = os.Environ()
	if err := wrapper.Listen(wrapCtx, cmd); err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "failed to start shell instance")
	}

	defer cancelFn()
	return cmd.Wait()
}
