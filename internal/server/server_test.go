package server_test

import (
	"fetcher/internal/server"
	"os"
	"testing"
	"time"
)

func TestServer_GracefulShutdown_Unit(t *testing.T) {
	interruptSigFunc := func(quit chan<- os.Signal) {
		time.Sleep(1 * time.Second)
		quit <- os.Interrupt
	}

	s := server.New()
	quit := make(chan os.Signal)
	done := make(chan struct{})

	go s.GracefulShutdown(quit, done)
	go interruptSigFunc(quit)

	err := s.ListenAndServe()
	if err != nil {
		t.Fatalf("expected to get nil err; got %v", err)
	}
}
