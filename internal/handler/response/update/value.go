package update

import (
	"net/http"
)

type ValueWriter struct{}

func NewValueWriter() *ValueWriter {
	return &ValueWriter{}
}

func (v *ValueWriter) WriteHeaders(resp http.ResponseWriter) {
	resp.Header().Set("Content-Type", "plain/text; charset=utf-8")
}

func (w *ValueWriter) WriteError(resp http.ResponseWriter, err string, code int) {
	http.Error(resp, err, code)
}

func (v *ValueWriter) WriteBody(resp http.ResponseWriter) error {
	resp.WriteHeader(http.StatusOK)
	return nil
}
