package storage

import "errors"

var (
	// ErrURLMappingNotFound is returned when a short URL doesn't exist
	ErrURLMappingNotFound = errors.New("url mapping not found")

	// ErrURLMappingExists is returned when trying to create a short URL that already exists
	ErrURLMappingExists = errors.New("url mapping already exists")

    // ErrOriginalURLExists is returned when trying to shorten a URL that's already shortened
    ErrOriginalURLExists = errors.New("original url already exists")
)
