package worker

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestProcessorProcessesRequestsInBackground(t *testing.T) {
	t.Parallel()

	var processed atomic.Int64
	p := NewProcessor(func(r Request) Result {
		processed.Add(1)
		return Result{
			ID:    r.ID,
			Value: r.Data + "_ok",
		}
	}, 16)
	defer p.Stop()

	const n = 50
	var wg sync.WaitGroup
	wg.Add(n)

	errCh := make(chan error, n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			res, err := p.Submit(ctx, Request{ID: i, Data: "v" + itoa(i)})
			if err != nil {
				errCh <- err
				return
			}
			want := "v" + itoa(i) + "_ok"
			if res.ID != i || res.Value != want {
				errCh <- &mismatchError{got: res.Value, want: want}
			}
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	if got := processed.Load(); got != n {
		t.Fatalf("expected %d processed requests, got %d", n, got)
	}
}

func TestProcessorStopCausesSubmitToFail(t *testing.T) {
	t.Parallel()

	p := NewProcessor(func(r Request) Result {
		time.Sleep(200 * time.Millisecond)
		return Result{ID: r.ID, Value: "done"}
	}, 1)

	p.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err := p.Submit(ctx, Request{ID: 1, Data: "x"})
	if err != ErrStopped {
		t.Fatalf("expected ErrStopped, got %v", err)
	}
}

func TestProcessorHonorsContextTimeout(t *testing.T) {
	t.Parallel()

	started := make(chan struct{})
	release := make(chan struct{})

	p := NewProcessor(func(r Request) Result {
		close(started)
		<-release // block until the test releases
		return Result{ID: r.ID, Value: "late"}
	}, 1)
	defer p.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Ensure the worker has started processing before we wait on ctx timeout.
	done := make(chan error, 1)
	go func() {
		_, err := p.Submit(ctx, Request{ID: 1, Data: "x"})
		done <- err
	}()

	select {
	case <-started:
	case <-time.After(1 * time.Second):
		t.Fatal("worker didn't start processing in time")
	}

	err := <-done
	if err == nil {
		t.Fatal("expected context deadline exceeded, got nil")
	}
	if ctx.Err() == nil {
		t.Fatal("expected ctx to be canceled/expired")
	}

	close(release)
}

type mismatchError struct {
	got  string
	want string
}

func (m *mismatchError) Error() string {
	return "mismatch"
}

func itoa(i int) string {
	// Local helper to avoid pulling in strconv in this tiny unit test.
	if i == 0 {
		return "0"
	}
	var buf [16]byte
	pos := len(buf)
	n := i
	if n < 0 {
		// Not needed for this test, but keeps helper correct.
		return "-" + itoa(-n)
	}
	for n > 0 {
		pos--
		buf[pos] = byte('0' + (n % 10))
		n /= 10
	}
	return string(buf[pos:])
}

