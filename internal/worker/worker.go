package worker

import (
	"bytes"
	"context"
	"fetcher/internal/storage"
	"fetcher/internal/storage/memory"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Config struct {
	URL         string
	Interval    time.Duration
	ReaderLimit int64
}

type Worker struct {
	cfg        Config
	logger     *logrus.Logger
	writer     storage.DataWriter
	cli        *http.Client
	terminated bool
}

func (w *Worker) Terminated() bool {
	return w.terminated
}

func (w *Worker) fetch(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, w.cfg.URL, nil)
	if err != nil {
		w.logger.Errorf("request with ctx creation failed: %v", err)
		return "", fmt.Errorf("request with ctx creation failed: %v", err)
	}
	res, err := w.cli.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request cli call to %s failed: %v", w.cfg.URL, err)
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request response from %s failed. Status code: %d", w.cfg.URL, res.StatusCode)
	}

	var buff bytes.Buffer
	lr := io.LimitReader(res.Body, w.cfg.ReaderLimit)
	_, err = buff.ReadFrom(lr)
	if err != nil {
		return "", fmt.Errorf("response body read failed; err: %v", err)
	}
	return buff.String(), nil
}

func (w *Worker) Run(ctx context.Context, done <-chan struct{}) {
	for {
		select {
		case <-time.After(w.cfg.Interval):
			UUID := uuid.New()
			t1 := time.Now()
			s, err := w.fetch(ctx)
			if err != nil {
				t2 := time.Now()
				diff := t2.Sub(t1)
				inner := w.writer.Write(ctx, UUID.String(), "NULL", diff)
				if inner != nil {
					w.logger.Errorf("fetch data write failed: %v", err)
				}
				continue
			}

			t2 := time.Now()
			diff := t2.Sub(t1)
			inner := w.writer.Write(ctx, UUID.String(), s, diff)
			if inner != nil {
				w.logger.Errorf("fetch data write failed: %v", err)
				continue
			}
		case <-done:
			w.terminated = true
			return
		case <-ctx.Done():
			w.logger.Infof("Fetch worker terminated; reason: %s", ctx.Err())
			w.terminated = true
			return
		}
	}
}

type Option func(w *Worker)

func WithLogger(l *logrus.Logger) Option {
	return func(w *Worker) {
		w.logger = l
	}
}

func WithWriter(dw storage.DataWriter) Option {
	return func(w *Worker) {
		w.writer = dw
	}
}

func WithConfig(outer Config) Option {
	return func(w *Worker) {
		if outer.URL != "" {
			w.cfg.URL = outer.URL
		}
		if outer.Interval > 0 {
			w.cfg.Interval = outer.Interval
		}
	}
}

func New(opts ...Option) *Worker {
	cli := http.Client{
		Transport: http.DefaultTransport,
		Timeout:   5 * time.Second,
	}

	w := Worker{
		cli:    &cli,
		logger: logrus.StandardLogger(),
		writer: memory.New(),
		cfg: Config{
			URL:         "https://httpbin.org/range/15",
			Interval:    5 * time.Second,
			ReaderLimit: 1 << 20, // eq. 1048576
		},
	}

	for _, o := range opts {
		o(&w)
	}
	return &w
}
