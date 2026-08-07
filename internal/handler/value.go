package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/params"
	models "github.com/EvgeniyAleksandrov/metrics-server/internal/model"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/repository"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/service"
)

type ValueMetricsProcessor interface {
	Get(valueParams params.Value) (*models.Metrics, error)
}

type ValueParamsParser interface {
	Parse(req *http.Request) (*params.Value, error)
}

type ValueResponseWriter interface {
	WriteHeaders(resp http.ResponseWriter)
	WriteError(resp http.ResponseWriter, err string, code int)
	WriteBody(resp http.ResponseWriter, metrics models.Metrics) error
}

type Value struct {
	processor      ValueMetricsProcessor
	paramsParser   ValueParamsParser
	responseWriter ValueResponseWriter
}

func NewValue(
	processor ValueMetricsProcessor,
	paramsParser ValueParamsParser,
	responseWriter ValueResponseWriter,
) *Value {
	return &Value{
		processor:      processor,
		paramsParser:   paramsParser,
		responseWriter: responseWriter,
	}
}

func (v *Value) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	v.responseWriter.WriteHeaders(resp)

	valueParams, err := v.paramsParser.Parse(req)
	if err != nil {
		switch {
		case v.OneOfErrorsIs(err, ErrInvalidParams, ErrUnsupportedMetricType):
			v.responseWriter.WriteError(resp, "Not Found", http.StatusNotFound)
		default:
			v.responseWriter.WriteError(resp, fmt.Sprintf("Parse request error: %s", err.Error()), http.StatusBadRequest)
		}

		return
	}

	metricValue, err := v.processor.Get(*valueParams)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUnsupportedProcessMethod) || errors.Is(err, repository.ErrNotFoundElement):
			v.responseWriter.WriteError(resp, "Not Found", http.StatusNotFound)
		default:
			v.responseWriter.WriteError(resp, "Not initialized error", http.StatusInternalServerError)
		}
		return
	}

	if err := v.responseWriter.WriteBody(resp, *metricValue); err != nil {
		v.responseWriter.WriteError(resp, "Write response error", http.StatusInternalServerError)
		return
	}
}

func (v *Value) OneOfErrorsIs(err error, targets ...error) bool {
	if len(targets) == 0 || err == nil {
		return false
	}

	for _, target := range targets {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}
