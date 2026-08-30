package middleware

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/hash"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/interfaces"
	"go.uber.org/zap"
)

const HashKey = "HashSHA256"

type Hasher interface {
	MakeHash(data []byte) []byte
}

type HashWriter struct {
	http.ResponseWriter
	key string
}

func NewHashWriter(resp http.ResponseWriter, key string) *HashWriter {
	return &HashWriter{
		ResponseWriter: resp,
		key:            key,
	}
}

func (w *HashWriter) Write(body []byte) (int, error) {
	w.ResponseWriter.Header().Set(HashKey, hex.EncodeToString(hash.NewSHA256(w.key).MakeHash(body)))
	return w.ResponseWriter.Write(body)
}

type Hashing struct {
	logger interfaces.Logger

	key string
}

func NewHashing(logger interfaces.Logger, key string) *Hashing {
	return &Hashing{
		logger: logger,
		key:    key,
	}
}

func (h *Hashing) Do(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		hashString := req.Header.Get(HashKey)

		hashResp := NewHashWriter(resp, h.key)

		if hashString == "" || h.key == "" {
			handler.ServeHTTP(hashResp, req)
			return
		}

		bodyData, err := io.ReadAll(req.Body)
		if err != nil {
			h.logger.Warn("Read request body error", zap.Error(err))

			http.Error(resp, "Failsed to read body", http.StatusBadRequest)
			return
		}

		req.Body = io.NopCloser(bytes.NewReader(bodyData))

		reqHashString := hex.EncodeToString(hash.NewSHA256(h.key).MakeHash(bodyData))

		fmt.Println("HASHING EQUAL: ", reqHashString, hashString)

		if reqHashString != hashString {
			h.logger.Warn("Not equal hash")

			http.Error(resp, "Wrong hash", http.StatusBadRequest)
			return
		}

		handler.ServeHTTP(hashResp, req)
	})
}
