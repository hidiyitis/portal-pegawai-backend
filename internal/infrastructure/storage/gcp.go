package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GCPStorage interface {
	UploadFile(ctx context.Context, bucketName string, file io.Reader, fileName string) (string, error)
}

type gcpStorageImpl struct {
	client *storage.Client
}

func NewGCPStorage(ctx context.Context, credentialsJSON string) (GCPStorage, error) {
	var client *storage.Client
	var err error

	if credentialsJSON != "" {
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(credentialsJSON))
	} else {
		client, err = storage.NewClient(ctx, option.WithCredentialsFile("storage-service-account.json"))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create GCP storage client: %w", err)
	}

	return &gcpStorageImpl{
		client: client,
	}, nil
}

func (g *gcpStorageImpl) UploadFile(
	ctx context.Context,
	bucketName string,
	file io.Reader,
	fileName string,
) (string, error) {
	uniqueFileName := generateUniqueFileName(fileName)

	bucket := g.client.Bucket(bucketName)
	obj := bucket.Object(uniqueFileName)
	w := obj.NewWriter(ctx)
	if _, err := io.Copy(w, file); err != nil {
		return "", fmt.Errorf("failed to write file to GCS: %w", err)
	}

	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to close GCS writer: %w", err)
	}

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, uniqueFileName), nil
}

func generateUniqueFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	name := originalName[:len(originalName)-len(ext)]
	return fmt.Sprintf("%s_%d%s", name, time.Now().UnixNano(), ext)
}
