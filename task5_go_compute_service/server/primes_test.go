package server

import (
	"context"
	"testing"
	"time"
)

func TestComputePrimeCount_Basic(t *testing.T) {
	t.Parallel()

	cases := []struct {
		limit int
		want  int
	}{
		{limit: 0, want: 0},
		{limit: 1, want: 0},
		{limit: 2, want: 1},
		{limit: 3, want: 2},
		{limit: 10, want: 4},   // 2,3,5,7
		{limit: 30, want: 10},  // 2..29 excluding composites
		{limit: 100, want: 25}, // known pi(100)=25
	}

	for _, tc := range cases {
		got, err := ComputePrimeCount(context.Background(), tc.limit)
		if err != nil {
			t.Fatalf("limit=%d err=%v", tc.limit, err)
		}
		if got != tc.want {
			t.Fatalf("limit=%d expected %d, got %d", tc.limit, tc.want, got)
		}
	}
}

func TestComputePrimeCount_ContextCancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	_, err := ComputePrimeCount(ctx, 5_000_000)
	if err == nil {
		t.Fatal("expected cancellation error, got nil")
	}
}
