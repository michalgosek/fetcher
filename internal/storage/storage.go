package storage

import (
	"context"
	"time"
)

type DataWriter interface {
	Write(ctx context.Context, ID, content string, duration time.Duration) error
}
