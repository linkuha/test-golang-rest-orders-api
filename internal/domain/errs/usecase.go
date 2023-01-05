package errs

import "errors"

var (
	InvalidPassword = errors.New("invalid password")
	LogicalError    = errors.New("logical error")
)
