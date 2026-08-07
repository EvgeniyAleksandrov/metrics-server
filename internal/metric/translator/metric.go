package translator

import (
	"errors"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
	models "github.com/EvgeniyAleksandrov/metrics-server/internal/model"
)

var ErrUnsupportedTranslation = errors.New("unsupported translation")

type Translation interface {
	Translate(m *metric.Metric) (*models.Metrics, error)
	Type() metric.ValueType
}
type Metric struct {
	translation map[metric.ValueType]Translation
}

func NewMetric(translations ...Translation) *Metric {
	t := &Metric{
		translation: make(map[metric.ValueType]Translation, len(translations)),
	}

	for _, translation := range translations {
		t.translation[translation.Type()] = translation
	}

	return t
}

func (t *Metric) Translate(m *metric.Metric) (*models.Metrics, error) {
	translation, ok := t.translation[m.Type()]
	if !ok {
		return nil, ErrUnsupportedTranslation
	}

	return translation.Translate(m)
}
