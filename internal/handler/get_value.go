package handler

import (
	"errors"
	"net/http"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/repository"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/service"
)

type GetMetricProcessor interface {
	Get(method string, name string) (string, error)
}

type GetValue struct {
	processor GetMetricProcessor
}

func NewGetValue(processor GetMetricProcessor) *GetValue {
	return &GetValue{
		processor: processor,
	}
}

func (v *GetValue) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	method := req.PathValue("method")
	name := req.PathValue("name")

	if name == "" || method == "" {
		http.Error(res, "Not Found", http.StatusNotFound)
		return
	}

	metricValue, err := v.processor.Get(method, name)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUnsupportedProcessMethod) || errors.Is(err, repository.ErrNotFoundElement):
			http.Error(res, "Not Found", http.StatusNotFound)
		default:
			http.Error(res, "Not initialized error", http.StatusInternalServerError)
		}
		return
	}

	res.WriteHeader(http.StatusOK)
	_, _ = res.Write([]byte(metricValue))
}
