package main

import (
	"fetcher/internal/config"
	"fetcher/internal/server"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
)

func main() {
	if err := execute(); err != nil {
		log.Fatal(err)
	}
}

func execute() error {
	var path string
	flag.StringVar(&path, "c", "./config.yml", "service config file path")
	flag.Parse()

	cfg, err := config.New(config.WithPath(path))
	if err != nil {
		return fmt.Errorf("config file read failed %v", err)
	}

	srv := server.New(server.WithConfig(cfg.Server))

	quit := make(chan os.Signal, 1)
	done := make(chan struct{})

	signal.Notify(quit, os.Interrupt)

	go srv.GracefulShutdown(quit, done)

	err = srv.ListenAndServe()
	if err != nil {
		quit <- os.Interrupt
		<-done
		return fmt.Errorf("listen and serve failed %v", err)
	}

	<-done
	log.Println("service has been gracefully shutdown ")
	return nil
}
