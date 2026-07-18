//go:generate mockgen -source=update.go -destination=update_mock_test.go -package=handler_test
package handler

import (
	"errors"
	"net/http"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/service"
)

type UpdateMetricProcessor interface {
	Update(method string, name string, value string) error
}

type Update struct {
	metricProcessor UpdateMetricProcessor
}

func NewUpdate(metricProcessor UpdateMetricProcessor) *Update {
	return &Update{
		metricProcessor: metricProcessor,
	}
}

func (u *Update) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(resp, "unsupported method", http.StatusMethodNotAllowed)
		return
	}

	method := req.PathValue("method")
	name := req.PathValue("name")
	value := req.PathValue("value")

	if method == "" || name == "" || value == "" {
		http.Error(resp, "Not Found", http.StatusNotFound)
		return
	}

	if err := u.metricProcessor.Update(method, name, value); err != nil {
		switch {
		case errors.Is(err, service.ErrUnsupportedProcessMethod):
			http.Error(resp, "Unsupported method", http.StatusBadRequest)

		case errors.Is(err, service.ErrInvalidValueFormat):
			http.Error(resp, "Invalid value format", http.StatusBadRequest)

		default:
			http.Error(resp, "Not initialized error", http.StatusInternalServerError)
		}

		return
	}

	u.writeHeaders(resp)
	resp.WriteHeader(http.StatusOK)
}

func (u *Update) writeHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "plain/text")
	w.Header().Set("charset", "utf-8")
}
