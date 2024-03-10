package model

import "errors"

var (
	ErrInvalidOption = errors.New("invalid option")

	ErrActonFailed  = errors.New("action failed")
	ErrAssertFailed = errors.New("assert failed")
)
