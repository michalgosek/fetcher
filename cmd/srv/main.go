package main

import (
	"context"
	"fetcher/internal/config"
	"fetcher/internal/server"
	"fetcher/internal/storage/memory"
	"fetcher/internal/worker"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
)

func main() {
	if err := execute(); err != nil {
		log.Fatal(err)
	}
}

// ./../../config.yml
func execute() error {
	var path string
	flag.StringVar(&path, "c", "./config.yml", "service config file path")
	flag.Parse()

	cfg, err := config.New(config.WithPath(path))
	if err != nil {
		return fmt.Errorf("config file read failed %v", err)
	}

	srv := server.New(server.WithConfig(cfg.Server))

	wr := memory.New()
	w := worker.New(worker.WithWriter(wr))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	done := make(chan struct{}, 1)
	signal.Notify(quit, os.Interrupt)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		srv.GracefulShutdown(quit, done)
		log.Println("GracefulShutdown go routine termiinated")
	}()

	go func() {
		defer wg.Done()
		w.Run(ctx, done)
		log.Println("Fetch worker go routine termiinated")
	}()

	err = srv.ListenAndServe()
	if err != nil {
		quit <- os.Interrupt
		return fmt.Errorf("listen and serve failed %v", err)
	}

	wg.Wait()
	log.Println("service has been gracefully shutdown ")
	return nil
}
