//go:generate mockgen -source=update.go -destination=update_mock_test.go -package=handler_test
package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/params"
)

type UpdateMetricProcessor interface {
	Update(ctx context.Context, updateParams params.Update) error
}

type UpdateParamsParser interface {
	Parse(req *http.Request) (*params.Update, error)
}

type UpdateResponseWriter interface {
	WriteHeaders(resp http.ResponseWriter)
	WriteError(resp http.ResponseWriter, err string, code int)
	WriteBody(resp http.ResponseWriter) error
}

type Update struct {
	metricProcessor UpdateMetricProcessor
	requestParser   UpdateParamsParser
	responseWriter  UpdateResponseWriter
}

func NewUpdate(metricProcessor UpdateMetricProcessor, requestParser UpdateParamsParser, responseWriter UpdateResponseWriter) *Update {
	return &Update{
		metricProcessor: metricProcessor,
		requestParser:   requestParser,
		responseWriter:  responseWriter,
	}
}

func (u *Update) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	u.responseWriter.WriteHeaders(resp)

	updateParams, err := u.requestParser.Parse(req)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnsupportedMetricType):
			u.responseWriter.WriteError(resp, "Unsupported method", http.StatusBadRequest)
		case errors.Is(err, ErrInvalidValueFormat):
			u.responseWriter.WriteError(resp, "Invalid value format", http.StatusBadRequest)
		case errors.Is(err, ErrInvalidParams):
			u.responseWriter.WriteError(resp, "Not found", http.StatusNotFound)
		default:
			u.responseWriter.WriteError(resp, "Request parsing error", http.StatusInternalServerError)
		}

		return
	}

	if updateParams == nil {
		u.responseWriter.WriteError(resp, "Unexpected params", http.StatusInternalServerError)
		return
	}

	updateContext, cancel := context.WithTimeout(req.Context(), time.Duration(5)*time.Second)
	defer cancel()

	if err := u.metricProcessor.Update(updateContext, *updateParams); err != nil {
		u.responseWriter.WriteError(resp, "Processing request failed", http.StatusInternalServerError)
		return
	}

	if err := u.responseWriter.WriteBody(resp); err != nil {
		u.responseWriter.WriteError(resp, "Write header error", http.StatusInternalServerError)
		return
	}
}
