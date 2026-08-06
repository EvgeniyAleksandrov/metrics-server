package value

import (
	"fmt"
	"net/http"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/types"
	models "github.com/EvgeniyAleksandrov/metrics-server/internal/model"
)

type ValueWriter struct{}

func NewValueWriter() *ValueWriter {
	return &ValueWriter{}
}

func (w *ValueWriter) WriteHeaders(resp http.ResponseWriter) {
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func (w *ValueWriter) WriteError(resp http.ResponseWriter, err string, code int) {
	http.Error(resp, err, code)
}

func (w *ValueWriter) WriteBody(resp http.ResponseWriter, metrics models.Metrics) error {
	var value string

	switch metrics.MType {
	case models.Gauge:
		if metrics.Value == nil {
			return handler.ErrEmptyMetricsData
		}
		value = types.Gauge(*metrics.Value).String()
	case models.Counter:
		if metrics.Delta == nil {
			return handler.ErrEmptyMetricsData
		}

		value = types.Counter(*metrics.Delta).String()
	default:
		return handler.ErrUnsupportedMetricType
	}

	if _, err := resp.Write([]byte(value)); err != nil {
		return fmt.Errorf("write response: %w", err)
	}

	resp.WriteHeader(http.StatusOK)

	return nil
}
