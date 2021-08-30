package worker_test

import (
	"context"
	"fetcher/internal/worker"
	"testing"
)

func TestWorker_RunCancellation_Unit(t *testing.T) {
	ctx := context.Background()
	w := worker.New()
	done := make(chan struct{})
	go w.Run(ctx, done)

	done <- struct{}{}
	if w.Terminated() != true {
		t.Fatalf("expected to get true for worker terminated; got %v", w.Terminated())
	}
}
