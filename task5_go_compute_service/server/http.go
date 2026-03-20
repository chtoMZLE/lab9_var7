package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type primeRequest struct {
	Limit int `json:"limit"`
}

type primeResponse struct {
	Limit      int `json:"limit"`
	PrimeCount int `json:"prime_count"`
}

var ErrBadRequest = errors.New("bad request")

func handlePrimes(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req primeRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		http.Error(w, ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	if req.Limit < 0 {
		http.Error(w, "limit must be >= 0", http.StatusBadRequest)
		return
	}

	count, err := ComputePrimeCount(ctx, req.Limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusRequestTimeout)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(primeResponse{
		Limit:      req.Limit,
		PrimeCount: count,
	})
}

// NewMux builds an HTTP handler mux.
func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/primes", func(w http.ResponseWriter, r *http.Request) {
		// r.Context() is canceled when client disconnects.
		handlePrimes(r.Context(), w, r)
	})
	return mux
}
