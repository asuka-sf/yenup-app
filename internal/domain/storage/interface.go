package storage

import (
	"context"
	"yenup/internal/domain/rate"
)

type Client interface {
	// Read a rate data from storage
	Read(ctx context.Context) ([]*rate.Rate, error)
	// Write a rate data to storage
	Write(ctx context.Context, rates []*rate.Rate) error
}
