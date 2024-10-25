package storage

import "context"

type Storage interface {
	SaveFile(ctx context.Context, path, filename, contentType string, fileBytes []byte) (string, error)
	DeleteFile(ctx context.Context, filename string) error
}
