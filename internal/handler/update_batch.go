//go:generate mockgen -source=update.go -destination=update_mock_test.go -package=handler_test
package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/params"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/interfaces"
	"go.uber.org/zap"
)

type UpdateBatchMetricProcessor interface {
	UpdateBatch(ctx context.Context, updateParams []params.Update) error
}

type UpdateBatchParamsParser interface {
	Parse(req *http.Request) ([]params.Update, error)
}

type UpdateBatchResponseWriter interface {
	WriteHeaders(resp http.ResponseWriter)
	WriteError(resp http.ResponseWriter, err string, code int)
	WriteBody(resp http.ResponseWriter) error
}

type UpdateBatch struct {
	metricProcessor UpdateBatchMetricProcessor
	requestParser   UpdateBatchParamsParser
	responseWriter  UpdateBatchResponseWriter
	logger          interfaces.Logger
}

func NewUpdateBatch(
	batchMetricProcessor UpdateBatchMetricProcessor,
	requestParser UpdateBatchParamsParser,
	responseWriter UpdateBatchResponseWriter,
	logger interfaces.Logger,
) *UpdateBatch {
	return &UpdateBatch{
		metricProcessor: batchMetricProcessor,
		requestParser:   requestParser,
		responseWriter:  responseWriter,
		logger:          logger,
	}
}

func (b *UpdateBatch) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	b.responseWriter.WriteHeaders(resp)

	updateParams, err := b.requestParser.Parse(req)
	if err != nil {
		b.logger.Warn("Parse request error", zap.Error(err))

		switch {
		case errors.Is(err, ErrUnsupportedMetricType):
			b.responseWriter.WriteError(resp, "Unsupported method", http.StatusBadRequest)
		case errors.Is(err, ErrInvalidValueFormat):
			b.responseWriter.WriteError(resp, "Invalid value format", http.StatusBadRequest)
		case errors.Is(err, ErrInvalidParams):
			b.responseWriter.WriteError(resp, "Not found", http.StatusNotFound)
		default:
			b.responseWriter.WriteError(resp, "Request parsing error", http.StatusInternalServerError)
		}

		return
	}

	if updateParams == nil {
		b.responseWriter.WriteError(resp, "Unexpected params", http.StatusInternalServerError)
		return
	}

	updateContext, cancel := context.WithTimeout(req.Context(), time.Duration(5)*time.Second)
	defer cancel()

	if err := b.metricProcessor.UpdateBatch(updateContext, updateParams); err != nil {
		b.logger.Warn("Processing metrics error", zap.Error(err))
		b.responseWriter.WriteError(resp, "Processing request failed", http.StatusInternalServerError)

		return
	}

	if err := b.responseWriter.WriteBody(resp); err != nil {
		b.logger.Warn("Write response error", zap.Error(err))
		b.responseWriter.WriteError(resp, "Write header error", http.StatusInternalServerError)

		return
	}
}
