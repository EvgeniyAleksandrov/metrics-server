package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string, fields ...zap.Field)
}

type loggingResponseWriter struct {
	http.ResponseWriter
	size   int
	status int
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.size += size

	return size, err
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *loggingResponseWriter) GetSize() int {
	return w.size
}

func (w *loggingResponseWriter) GetStatus() int {
	return w.status
}

type Logging struct {
	logger Logger
}

func NewLogging(logger Logger) *Logging {
	return &Logging{
		logger: logger,
	}
}

func (l *Logging) Do(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		start := time.Now()

		lResp := &loggingResponseWriter{
			ResponseWriter: resp,
		}

		handler.ServeHTTP(lResp, req)

		l.logger.Info(
			"Request and Response data",
			zap.String("uri", req.RequestURI),
			zap.String("method", req.Method),
			zap.Int("responseCode", lResp.GetStatus()),
			zap.Int("responseSize", lResp.GetSize()),
			zap.Duration("responseTime", time.Now().Sub(start)),
		)
	})
}
