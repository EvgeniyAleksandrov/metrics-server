package middlware

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

func WithLogging(handlerFunc http.HandlerFunc, logger Logger) http.HandlerFunc {
	l := func(resp http.ResponseWriter, req *http.Request) {
		start := time.Now()

		lResp := &loggingResponseWriter{
			ResponseWriter: resp,
		}

		handlerFunc(lResp, req)

		logger.Info(
			"Resuest and Response data",
			zap.String("uri", req.RequestURI),
			zap.String("method", req.Method),
			zap.Int("responseCode", lResp.GetStatus()),
			zap.Int("responseSize", lResp.GetSize()),
			zap.Duration("responseTime", time.Now().Sub(start)),
		)
	}

	return l
}
