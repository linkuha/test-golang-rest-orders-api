package v1

import "github.com/pkg/errors"

var (
	MalformedRequest = errors.New("malformed request body")
	ServiceError     = errors.New("server error, try later")
)
