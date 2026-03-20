package tcpsrv

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestTCPServer_UppercasePerLine(t *testing.T) {
	t.Parallel()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	srv := NewServer(l, func(ctx context.Context, msg string) (string, error) {
		return strings.ToUpper(msg), nil
	})

	go func() { _ = srv.Serve() }()
	t.Cleanup(func() { srv.Stop() })

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	// Send multiple lines on a single connection.
	_, _ = fmt.Fprint(conn, "hello\nworld\n")

	r := bufio.NewReader(conn)
	line1, err := r.ReadString('\n')
	if err != nil {
		t.Fatalf("read1: %v", err)
	}
	line2, err := r.ReadString('\n')
	if err != nil {
		t.Fatalf("read2: %v", err)
	}

	line1 = strings.TrimSpace(line1)
	line2 = strings.TrimSpace(line2)

	if line1 != "HELLO" {
		t.Fatalf("expected HELLO, got %q", line1)
	}
	if line2 != "WORLD" {
		t.Fatalf("expected WORLD, got %q", line2)
	}
}

func TestTCPServer_HandlesConcurrentClients(t *testing.T) {
	t.Parallel()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	srv := NewServer(l, func(ctx context.Context, msg string) (string, error) {
		return strings.ToUpper(msg), nil
	})

	go func() { _ = srv.Serve() }()
	t.Cleanup(func() { srv.Stop() })

	const clients = 25
	const perClient = 5

	var wg sync.WaitGroup
	wg.Add(clients)

	errCh := make(chan error, clients)
	for c := 0; c < clients; c++ {
		c := c
		go func() {
			defer wg.Done()

			conn, err := net.Dial("tcp", l.Addr().String())
			if err != nil {
				errCh <- fmt.Errorf("dial: %w", err)
				return
			}
			defer conn.Close()

			msg := make([]string, 0, perClient)
			for i := 0; i < perClient; i++ {
				msg = append(msg, fmt.Sprintf("c%d-%d", c, i))
			}

			for _, m := range msg {
				_, _ = fmt.Fprintf(conn, "%s\n", m)
			}

			r := bufio.NewReader(conn)
			for _, m := range msg {
				got, err := r.ReadString('\n')
				if err != nil {
					errCh <- fmt.Errorf("read: %w", err)
					return
				}
				got = strings.TrimSpace(got)
				want := strings.ToUpper(m)
				if got != want {
					errCh <- fmt.Errorf("expected %q, got %q", want, got)
					return
				}
			}
		}()
	}

	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestTCPServer_StopClosesListener(t *testing.T) {
	t.Parallel()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	srv := NewServer(l, func(ctx context.Context, msg string) (string, error) {
		return msg, nil
	})
	go func() { _ = srv.Serve() }()

	// Stop quickly, then ensure new dials fail.
	time.Sleep(50 * time.Millisecond)
	srv.Stop()

	dialer := &net.Dialer{Timeout: 200 * time.Millisecond}
	conn, err := dialer.Dial("tcp", l.Addr().String())
	if err == nil {
		conn.Close()
		t.Fatal("expected dial error after Stop(), got nil")
	}
}
