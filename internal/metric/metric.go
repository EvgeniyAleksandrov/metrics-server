package metric

type Value interface {
	Type() ValueType
}

type Metric struct {
	name  Name
	value Value
}

func NewMetric(name Name, value Value) *Metric {
	return &Metric{
		name:  name,
		value: value,
	}
}

func (m *Metric) Name() Name {
	return m.name
}

func (m *Metric) Value() any {
	return m.value
}

func (m *Metric) Type() ValueType {
	return m.value.Type()
}
