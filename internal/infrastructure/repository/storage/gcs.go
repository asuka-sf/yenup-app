package storage

import (
	"context"
	"encoding/json"
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
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var rates []*rate.Rate
	err = json.Unmarshal(data, &rates)
	if err != nil {
		return nil, err
	}

	return rates, err

}

// Write serializes rate data and saves it to GCS.
// func (g *GCSClient) Write(rates []*rate.Rate) error {}
