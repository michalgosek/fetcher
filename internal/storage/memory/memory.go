package memory

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

type Memory struct {
	lookup sync.Map
}

type Record struct {
	UUID     string
	Content  string
	Duration time.Duration
}

func (m *Memory) Write(ctx context.Context, UUID, c string, d time.Duration) error {
	_, ok := m.lookup.Load(UUID)
	if ok {
		return ErrUUIDNotUnique
	}

	r := Record{
		UUID:     UUID,
		Content:  c,
		Duration: d,
	}
	m.lookup.Store(UUID, r)

	log.Println("Inserted:")
	log.Printf("%+v", r)
	return nil
}

func New() *Memory {
	return &Memory{}
}

var ErrUUIDNotUnique = errors.New("provided UUID already exists; write op failed")
