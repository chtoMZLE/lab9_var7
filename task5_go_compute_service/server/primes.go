package server

import (
	"context"
)

// ComputePrimeCount returns the number of primes <= limit.
//
// It uses a sieve of Eratosthenes and periodically checks ctx for cancellation.
func ComputePrimeCount(ctx context.Context, limit int) (int, error) {
	if limit < 2 {
		return 0, nil
	}
	// ctx check frequency: every chunk of candidate updates.
	const checkEvery = 4096

	// sieve[i] == true means "i is prime candidate".
	sieve := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		sieve[i] = true
	}

	for p := 2; p*p <= limit; p++ {
		if !sieve[p] {
			continue
		}
		// Avoid overflow in p*p.
		start := p * p
		for m := start; m <= limit; m += p {
			sieve[m] = false
			// Periodically check for cancellation.
			if m%checkEvery == 0 {
				select {
				case <-ctx.Done():
					return 0, ctx.Err()
				default:
				}
			}
		}
	}

	count := 0
	for i := 2; i <= limit; i++ {
		if sieve[i] {
			count++
		}
	}
	return count, nil
}
