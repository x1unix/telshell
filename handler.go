package telshell

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"os"
	"os/exec"
)

const bufferSize = 4

type Handler interface {
	Handle(ctx context.Context, rw io.ReadWriter) error
}

type ShellHandler struct {
	shellPath string
	shellArgs []string
	log *zap.SugaredLogger
}

func NewShellHandler(shell string, args ...string) ShellHandler {
	return ShellHandler{
		shellPath: shell,
		shellArgs: args,
		log: zap.S().Named("shell"),
	}
}

func (s ShellHandler) Handle(ctx context.Context, rw io.ReadWriter) error {
	fmt.Fprintf(rw, "Current shell is %q\n", s.shellPath)

	wrapCtx, cancelFn := context.WithCancel(ctx)
	wrapper := NewTerminalWrapper(s.log, rw, bufferSize)
	cmd := exec.CommandContext(ctx, s.shellPath, "-i")
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
