package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/interfaces"
	"go.uber.org/zap"
)

const (
	ContentTypeTextHtml        = "text/html"
	ContentTypeApplicationJSON = "application/json"
)

type Compression struct {
	logger interfaces.Logger
}

func NewCompression(logger interfaces.Logger) *Compression {
	return &Compression{
		logger: logger,
	}
}

func (c *Compression) Do(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		oldResponseWriter := resp

		if c.isAcceptEncoding(req) {
			gzipWriter, err := gzip.NewWriterLevel(resp, gzip.BestSpeed)
			if err != nil {
				c.logger.Warn("Invalid compression data", zap.Error(err))
				http.Error(resp, "Invalid encoding data", http.StatusInternalServerError)
				return
			}

			defer gzipWriter.Close()

			oldResponseWriter = newCompressWriter(resp, gzipWriter)
		}

		if c.isContentEncoded(req) {
			cr, err := newCompressReader(req.Body)
			if err != nil {
				c.logger.Warn("Create compression reader", zap.Error(err))
				http.Error(resp, "Create compression error", http.StatusInternalServerError)
				return
			}
			defer cr.Close()
			req.Body = cr
		}

		handler.ServeHTTP(
			newResponseSupportedTypeWriter(
				oldResponseWriter,
				resp,
				ContentTypeTextHtml,
				ContentTypeApplicationJSON,
			),
			req,
		)
	})
}

func (c *Compression) isAcceptEncoding(req *http.Request) bool {
	return strings.Contains(req.Header.Get("Accept-Encoding"), "gzip")
}

func (c *Compression) isContentEncoded(req *http.Request) bool {
	return strings.Contains(req.Header.Get("Content-Encoding"), "gzip")
}

type compressWriter struct {
	http.ResponseWriter
	writer io.Writer
}

func newCompressWriter(responseWriter http.ResponseWriter, writer io.Writer) *compressWriter {
	return &compressWriter{
		ResponseWriter: responseWriter,
		writer:         writer,
	}
}

func (w *compressWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

type compressReader struct {
	reader     io.ReadCloser
	gzipReader *gzip.Reader
}

func newCompressReader(reader io.ReadCloser) (*compressReader, error) {
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}

	return &compressReader{reader: reader, gzipReader: gzipReader}, nil
}

func (r compressReader) Read(p []byte) (n int, err error) {
	return r.gzipReader.Read(p)
}

func (r *compressReader) Close() error {
	if err := r.reader.Close(); err != nil {
		return err
	}

	return r.gzipReader.Close()
}

type responseSupportedTypeWriter struct {
	http.ResponseWriter
	rawWriter http.ResponseWriter

	compressedContentTypes []string
}

func newResponseSupportedTypeWriter(
	responseWriter http.ResponseWriter,
	rawWriter http.ResponseWriter,
	compressedContentTypes ...string,
) *responseSupportedTypeWriter {
	return &responseSupportedTypeWriter{
		ResponseWriter:         responseWriter,
		rawWriter:              rawWriter,
		compressedContentTypes: compressedContentTypes,
	}
}

func (w *responseSupportedTypeWriter) Write(b []byte) (int, error) {
	if w.isSupportedContentType() {
		w.rawWriter.Header().Set("Content-Encoding", "gzip")
		return w.ResponseWriter.Write(b)
	}

	return w.rawWriter.Write(b)
}

func (w *responseSupportedTypeWriter) WriteHeader(statusCode int) {
	if w.isSupportedContentType() {
		w.rawWriter.Header().Set("Content-Encoding", "gzip")
	}

	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseSupportedTypeWriter) isSupportedContentType() bool {
	contentType := w.Header().Get("Content-Type")
	for _, supportedType := range w.compressedContentTypes {
		if strings.Contains(contentType, supportedType) {
			return true
		}
	}

	return false
}
