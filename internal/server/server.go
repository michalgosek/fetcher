package server

import (
	"log"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Server struct {
	logger *logrus.Logger
	http   *http.Server
}

func (s *Server) ListenAndServe() error {
	log.Printf("Server starts listen on port addr: %s\n", s.http.Addr)
	return s.http.ListenAndServe()
}

func New(l *logrus.Logger, h http.Handler) *Server {
	s := Server{
		logger: l,
		http: &http.Server{
			Addr:           "localhost:8080",
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
			Handler:        h,
		},
	}
	return &s
}
