package storage

import (
	"context"
	rate "yenup/internal/domain/rate"
)

type Client interface {
	// Read a JSON file
	Read(ctx context.Context) ([]*rate.Rate, error)
	// Write a JSON file
	// Write(path string, rates rate.Rate) error
}
