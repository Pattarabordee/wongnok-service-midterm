package global

import "errors"

var (
	ErrForbidden error = errors.New("forbidden")
)
