package errs

import "errors"

var (
	ErrName = errors.New("invalid name")
	ErrSign = errors.New("invalid sign")
	ErrPrice = errors.New("invalid price")
)