package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"cloud.google.com/go/storage"

	rate "yenup/internal/domain/rate"
)

// GCSClient wraps a GCS bucket handle and provides read/write operations for rate data.
type GCSClient struct {
	bucket *storage.BucketHandle
	object string
}

// NewGCSClient creates a new GCSClient with the specified bucket and object.
func NewGCSClient(client *storage.Client, bucketName, objectName string) *GCSClient {
	return &GCSClient{
		bucket: client.Bucket(bucketName),
		object: objectName,
	}
}

// Read fetches rate data from GCS and deserializes it into a slice of Rate.
func (g *GCSClient) Read(ctx context.Context) ([]*rate.Rate, error) {

	reader, err := g.bucket.Object(g.object).NewReader(ctx)
	// if the JSON file doesn't exist, return an empty slice
	if errors.Is(err, storage.ErrObjectNotExist) {
		return []*rate.Rate{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open reader: %w", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read GCS: %w", err)
	}

	var rates []*rate.Rate
	err = json.Unmarshal(data, &rates)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal rates: %w", err)
	}

	return rates, nil
}

// Write serializes rate data and saves it to GCS.
func (g *GCSClient) Write(ctx context.Context, rates []*rate.Rate) error {
	writer := g.bucket.Object(g.object).NewWriter(ctx)
	writer.ContentType = "application/json"

	rateJSON, err := json.Marshal(rates)
	if err != nil {
		return fmt.Errorf("failed to marshal rates: %w", err)
	}
	if _, err := writer.Write(rateJSON); err != nil {
		return fmt.Errorf("failed to write json: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close GCS writer: %w", err)
	}

	return nil
}
