package telshell

import (
	"bytes"
	"context"
	"io"
	"os/exec"

	"github.com/x1unix/telshell/internal/helpers"
	"go.uber.org/zap"
)

type IOParams struct {
	BufferSize         uint
	ReplaceLineEndings bool
}

type TerminalWrapper struct {
	IOParams
	client io.ReadWriter
	log    *zap.SugaredLogger
}

func NewTerminalWrapper(log *zap.SugaredLogger, client io.ReadWriter, params IOParams) TerminalWrapper {
	return TerminalWrapper{client: client, log: log, IOParams: params}
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

		buff := make([]byte, w.BufferSize)
		_, _ = r.Read(buff)

		if w.ReplaceLineEndings {
			// Replace RF with CRLF line endings
			buff = bytes.ReplaceAll(buff, []byte("\n"), []byte{CL, RF})
		}

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
		arr := make([]byte, w.BufferSize)
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
		case CL, NulChar:
			continue
		default:
		}

		filtered = append(filtered, b)
	}

	return filtered
}
