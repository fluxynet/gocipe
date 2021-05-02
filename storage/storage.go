package storage

import "context"

// Storage manages file access
type Storage interface {
	// Store a file
	Store(ctx context.Context, filename string, contents []byte) error

	// Delete a file or path
	Delete(ctx context.Context, path string) error
}
