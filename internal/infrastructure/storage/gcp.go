package storage

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"io"
	"os"
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

func NewGCPStorage(ctx context.Context) (GCPStorage, error) {

	client, err := storage.NewClient(ctx)
	if err == nil {
		return &gcpStorageImpl{client: client}, nil
	}

	const credsFile = "storage-service-account.json"
	if _, err := os.Stat(credsFile); err == nil {
		client, err := storage.NewClient(ctx, option.WithCredentialsFile(credsFile))
		if err != nil {
			return nil, fmt.Errorf("failed to create GCP storage client with file: %w", err)
		}
		return &gcpStorageImpl{client: client}, nil
	}

	if jsonCreds := os.Getenv("GOOGLE_CREDENTIALS_JSON"); jsonCreds != "" {
		creds, err := google.CredentialsFromJSON(ctx, []byte(jsonCreds), storage.ScopeReadWrite)
		if err != nil {
			return nil, fmt.Errorf("failed to parse credentials from env: %w", err)
		}
		client, err := storage.NewClient(ctx, option.WithCredentials(creds))
		if err != nil {
			return nil, fmt.Errorf("failed to create GCP storage client with env creds: %w", err)
		}
		return &gcpStorageImpl{client: client}, nil
	}

	return nil, fmt.Errorf("no valid Google Cloud credentials found")
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
