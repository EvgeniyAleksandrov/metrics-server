package update

import (
	"fmt"
	"net/http"
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

func (w *JSONWriter) WriteBody(resp http.ResponseWriter) error {
	resp.WriteHeader(http.StatusOK)

	fmt.Fprintln(resp, "{}")

	return nil
}
