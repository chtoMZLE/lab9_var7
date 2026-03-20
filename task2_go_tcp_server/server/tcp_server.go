package tcpsrv

import (
	"bufio"
	"context"
	"net"
	"strings"
	"sync"
)

// Handler converts an inbound line into an outbound line.
type Handler func(ctx context.Context, msg string) (string, error)

type Server struct {
	l       net.Listener
	handler Handler

	stopOnce sync.Once
	stopCh   chan struct{}

	wg sync.WaitGroup
}

func NewServer(l net.Listener, handler Handler) *Server {
	return &Server{
		l:       l,
		handler: handler,
		stopCh:  make(chan struct{}),
	}
}

// Serve starts accepting connections and processing them concurrently.
// It returns nil when Stop() is called.
func (s *Server) Serve() error {
	for {
		conn, err := s.l.Accept()
		if err != nil {
			select {
			case <-s.stopCh:
				return nil
			default:
				return err
			}
		}

		s.wg.Add(1)
		go func(c net.Conn) {
			defer s.wg.Done()
			s.handleConn(c)
		}(conn)
	}
}

func (s *Server) Stop() {
	s.stopOnce.Do(func() {
		close(s.stopCh)
		_ = s.l.Close()
	})
	s.wg.Wait()
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewScanner(conn)
	// Increase max token size to allow longer lines in tests/manual usage.
	reader.Buffer(make([]byte, 0, 4096), 1024*1024)
	writer := bufio.NewWriter(conn)

	ctx := context.Background()
	for reader.Scan() {
		line := strings.TrimRight(reader.Text(), "\r")
		out, err := s.handler(ctx, line)
		if err != nil {
			out = ""
		}

		// Line-based protocol: one response per input line.
		_, _ = writer.WriteString(out)
		_, _ = writer.WriteString("\n")
		_ = writer.Flush()
	}
}
