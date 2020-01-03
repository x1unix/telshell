package telshell

import (
	"context"
	"github.com/x1unix/telshell/internal/helpers"
	"go.uber.org/zap"
	"io"
	"os/exec"
)

type TerminalWrapper struct {
	buffSize int
	client   io.ReadWriter
	log      *zap.SugaredLogger
}

func NewTerminalWrapper(log *zap.SugaredLogger, client io.ReadWriter, buffSize int) TerminalWrapper {
	return TerminalWrapper{client: client, log: log, buffSize: buffSize}
}

func (w TerminalWrapper) Listen(ctx context.Context, cmd *exec.Cmd) error {
	return w.listenHost(ctx, cmd)
}

func (w TerminalWrapper) listenHost(ctx context.Context, cmd *exec.Cmd) error {
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	go w.readFromHost(ctx, stdout)
	go w.readFromHost(ctx, stderr)
	go w.writeToHost(ctx, stdin)

	return nil
}

func (w TerminalWrapper) readFromHost(ctx context.Context, r io.Reader) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		buff := make([]byte, w.buffSize)
		_, _ = r.Read(buff)
		w.client.Write(buff)
	}
}

func (w TerminalWrapper) writeToHost(ctx context.Context, dest io.Writer) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		arr := make([]byte, w.buffSize)
		_, err := w.client.Read(arr)
		if err == io.EOF || helpers.IsErrClosing(err) {
			return
		}
		if err != nil {
			w.log.Error(err)
			continue
		}

		arr = w.filterChars(arr)
		dest.Write(arr)
	}
}

func (w TerminalWrapper) filterChars(msg []byte) []byte {
	filtered := make([]byte, 0, len(msg))
	for _, b := range msg {
		switch b {
		case CR, NulChar:
			continue
		default:
		}

		filtered = append(filtered, b)
	}

	return filtered
}
