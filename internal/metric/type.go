package metric

type ValueType string

const (
	Gauge   ValueType = "gauge"
	Counter ValueType = "counter"
)
