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

type Hasher interface {
	MakeHash(data []byte) []byte
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
		hashString := req.Header.Get("HashSHA256")

		if hashString == "" || h.key == "" {
			handler.ServeHTTP(resp, req)
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

		handler.ServeHTTP(resp, req)
	})
}
