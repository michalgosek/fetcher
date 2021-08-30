package server_test

import (
	"fetcher/internal/server"
	"os"
	"sync"
	"testing"
	"time"
)

func TestServer_GracefulShutdown_Unit(t *testing.T) {
	var wg sync.WaitGroup
	interruptSigFunc := func(quit chan<- os.Signal) {
		defer wg.Done()
		time.Sleep(1 * time.Second)
		quit <- os.Interrupt
	}

	s := server.New()
	quit := make(chan os.Signal)
	done := make(chan struct{})

	wg.Add(1)
	go s.GracefulShutdown(quit, done)
	go interruptSigFunc(quit)

	wg.Wait()
	err := s.ListenAndServe()
	if err != nil {
		t.Fatalf("expected to get nil err; got %v", err)
	}
}
