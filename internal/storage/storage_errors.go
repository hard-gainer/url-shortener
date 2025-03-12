package storage

import "errors"

var (
	// ErrURLMappingNotFound is returned when a short URL doesn't exist
	ErrURLMappingNotFound = errors.New("url mapping not found")

	// ErrURLMappingExists is returned when trying to create a short URL that already exists
	ErrURLMappingExists = errors.New("url mapping already exists")
)
