package worker

import (
	"context"
	"errors"
	"sync"
)

// Request represents a unit of work to be processed asynchronously.
type Request struct {
	ID   int
	Data string
}

// Result represents the outcome of processing a Request.
type Result struct {
	ID    int
	Value string
}

var ErrStopped = errors.New("worker stopped")

type pendingRequest struct {
	req   Request
	reply chan Result
}

// Processor processes requests in a dedicated background goroutine.
//
// It guarantees that submitted requests are handled in FIFO order (as observed by the internal channel),
// and Submit either returns a result, ctx error, or ErrStopped.
type Processor struct {
	processFn func(Request) Result

	reqCh   chan pendingRequest
	stopCh  chan struct{}
	stopOnce sync.Once

	startedOnce sync.Once
}

// NewProcessor creates a processor and starts its background goroutine immediately.
func NewProcessor(processFn func(Request) Result, buffer int) *Processor {
	if buffer <= 0 {
		buffer = 1
	}
	p := &Processor{
		processFn: processFn,
		reqCh:     make(chan pendingRequest, buffer),
		stopCh:    make(chan struct{}),
	}
	p.startedOnce.Do(func() {
		// Start the background loop without blocking NewProcessor.
		go p.run()
	})
	return p
}

// Stop stops the background goroutine.
func (p *Processor) Stop() {
	p.stopOnce.Do(func() {
		close(p.stopCh)
	})
}

// Submit sends a request for background processing.
func (p *Processor) Submit(ctx context.Context, r Request) (Result, error) {
	pr := pendingRequest{
		req:   r,
		reply: make(chan Result, 1),
	}

	select {
	case <-p.stopCh:
		return Result{}, ErrStopped
	case <-ctx.Done():
		return Result{}, ctx.Err()
	case p.reqCh <- pr:
	}

	select {
	case <-p.stopCh:
		return Result{}, ErrStopped
	case <-ctx.Done():
		return Result{}, ctx.Err()
	case res := <-pr.reply:
		return res, nil
	}
}

func (p *Processor) run() {
	for {
		select {
		case <-p.stopCh:
			return
		case pr := <-p.reqCh:
			// pr.reply is buffered; Submit won't block even if the receiver timed out/canceled.
			res := p.processFn(pr.req)
			select {
			case pr.reply <- res:
			default:
				// Reply buffer is full; receiver already gave up (timeout/cancel).
			}
		}
	}
}

