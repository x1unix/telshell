package telshell

import (
	"context"
	"fmt"
	"github.com/x1unix/telshell/internal/helpers"
	"net"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	ErrServerStarted    = errors.New("server is already started")
	ErrServerNotStarted = errors.New("server not started")
)

type Server struct {
	log     *zap.SugaredLogger
	handler Handler
	ln      net.Listener

	running bool
	runLock sync.RWMutex

	ctx      context.Context
	cancelFn context.CancelFunc
}

func NewServer(handler Handler) *Server {
	return &Server{
		handler: handler,
		log:     zap.S().Named("server"),
	}
}

func (s *Server) SetLogger(l *zap.Logger) {
	s.log = l.Sugar()
}

func (s *Server) Start(ctx context.Context, addr string) (err error) {
	if err := s.checkRunState(true); err != nil {
		return err
	}

	var lc net.ListenConfig
	s.ctx, s.cancelFn = context.WithCancel(ctx)
	s.ln, err = lc.Listen(s.ctx, "tcp", addr)
	if err != nil {
		s.running = false
		return err
	}

	s.markRunState(true)
	s.log.Infof("listening on %q...", addr)

	return s.listen()
}

func (s *Server) checkRunState(isRunning bool) error {
	s.runLock.RLock()
	defer s.runLock.RUnlock()

	if isRunning == s.running {
		if isRunning {
			return ErrServerStarted
		}

		return ErrServerNotStarted
	}

	return nil
}

func (s *Server) markRunState(isRunning bool) {
	s.runLock.Lock()
	defer s.runLock.Unlock()
	s.running = isRunning
}

func (s *Server) listen() error {
	go func() {
		<-s.ctx.Done()
		s.log.Info("context is dead, closing server")
		if err := s.shutdown(); err != nil {
			s.log.Error(err)
		}
	}()

	for {
		select {
		case <-s.ctx.Done():
			s.log.Debug("context is dead, listener stop")
			return nil
		default:
		}

		conn, err := s.ln.Accept()
		if err != nil {
			if helpers.IsErrClosing(err) {
				return nil
			}

			if e := s.shutdown(); e != nil {
				s.log.Warn(e)
			}
			return err
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			s.log.Errorf("recovered from panic: %s", r)
		}
	}()

	s.log.Debugf("%q: received new connection", conn.RemoteAddr().String())
	fmt.Fprintln(conn, "Wellcome to TelShell :)\r\n")

	if err := s.handler.Handle(s.ctx, conn); err != nil {
		s.log.Errorf("handler returned an error: %s", err)
		fmt.Fprintf(conn, "ERROR:\t%s\r\n", err.Error())
	}

	// Close connection if wasn't closed
	_ = conn.Close()
	s.log.Debugf("%q: connection closed", conn.RemoteAddr().String())
}

func (s *Server) shutdown() error {
	s.markRunState(false)
	if err := s.ln.Close(); err != nil {
		return errors.Wrap(err, "failed to close TCP server")
	}

	return nil
}

func (s *Server) Stop() error {
	if err := s.checkRunState(false); err != nil {
		return err
	}

	s.cancelFn()
	return nil
}
