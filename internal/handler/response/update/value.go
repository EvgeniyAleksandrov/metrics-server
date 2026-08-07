package update

import (
	"net/http"
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

func (w *ValueWriter) WriteBody(resp http.ResponseWriter) error {
	resp.WriteHeader(http.StatusOK)
	return nil
}
