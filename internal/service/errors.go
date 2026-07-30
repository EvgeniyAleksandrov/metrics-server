package service

import "errors"

var (
	ErrUnsupportedProcessMethod = errors.New("unsupported process method")
	ErrInvalidValueFormat       = errors.New("invalid value format")
)
