package service

import "errors"

var (
	ErrUnsupportedProcessMethod = errors.New("unsupported process method")
	ErrInvalidValueFormat       = errors.New("invalid value format")
	ErrEmptyName                = errors.New("empty name")
	ErrEmptyValue               = errors.New("empty value")
)
