package values

type Gauge struct {
	Name  string
	Value float64
}

type Counter struct {
	Name  string
	Delta int64
}
