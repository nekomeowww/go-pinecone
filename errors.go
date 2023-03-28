package pinecone

import "errors"

var (
	// ErrRequestFailed is returned and wrapped when a request to the Pinecone API fails.
	ErrRequestFailed = errors.New("request failed")
	// ErrIndexNotFound is returned when an index is not found.
	ErrIndexNotFound = errors.New("index not found")
	// ErrInvalidParams is returned when an invalid parameter is passed to a function.
	ErrInvalidParams = errors.New("invalid params")
)
