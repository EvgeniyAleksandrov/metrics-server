package value

import (
	"encoding/json"
	"fmt"
	"net/http"

	models "github.com/EvgeniyAleksandrov/metrics-server/internal/model"
)

type JSONWriter struct{}

func NewJSONWriter() *JSONWriter {
	return &JSONWriter{}
}

func (w *JSONWriter) WriteHeaders(resp http.ResponseWriter) {
	resp.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func (w *JSONWriter) WriteError(resp http.ResponseWriter, error string, code int) {
	resp.Header().Del("Content-Length")

	w.WriteHeaders(resp)
	resp.WriteHeader(code)

	jsonError := fmt.Sprintf(`{"error": "%s"}`, error)

	fmt.Fprintln(resp, jsonError)
}

func (w *JSONWriter) WriteBody(resp http.ResponseWriter, metrics models.Metrics) error {
	metricData, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("marshling metric to json data: %w", err)
	}

	if _, err = resp.Write(metricData); err != nil {
		return fmt.Errorf("write response data: %w", err)
	}

	resp.WriteHeader(http.StatusOK)

	return nil
}
