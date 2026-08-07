package handler

import "errors"

var (
	ErrUnsupportedMetricType = errors.New("unsupported metric type")
	ErrInvalidValueFormat    = errors.New("invalid value format")
	ErrInvalidParams         = errors.New("invalid request params")

	ErrEmptyMetricsData = errors.New("empty metrics data")
)
