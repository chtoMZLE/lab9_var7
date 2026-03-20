package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandlePrimes_InvalidMethod(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/primes", nil)
	w := httptest.NewRecorder()
	handlePrimes(context.Background(), w, req)

	if w.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Result().StatusCode)
	}
}

func TestHandlePrimes_BadJSON(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/primes", bytes.NewBufferString("{"))
	w := httptest.NewRecorder()
	handlePrimes(context.Background(), w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Result().StatusCode)
	}
}

func TestHandlePrimes_OK(t *testing.T) {
	t.Parallel()

	body, _ := json.Marshal(primeRequest{Limit: 30})
	req := httptest.NewRequest(http.MethodPost, "/primes", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handlePrimes(context.Background(), w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Result().StatusCode)
	}

	var resp primeResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Limit != 30 || resp.PrimeCount != 10 {
		// primes <= 30 are: 2,3,5,7,11,13,17,19,23,29 => 10
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestHandlePrimes_RespectsContextCancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	body, _ := json.Marshal(primeRequest{Limit: 1_000_000})
	req := httptest.NewRequest(http.MethodPost, "/primes", bytes.NewReader(body))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handlePrimes(ctx, w, req)

	// Compute may start but should fail fast with ctx cancellation.
	if w.Result().StatusCode != http.StatusRequestTimeout && w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected timeout/bad request, got %d", w.Result().StatusCode)
	}
}
