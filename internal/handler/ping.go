package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/interfaces"
	"go.uber.org/zap"
)

type Ping struct {
	storage interfaces.Storage
	logger  interfaces.Logger
}

func NewPing(storage interfaces.Storage, logger interfaces.Logger) *Ping {
	return &Ping{
		storage: storage,
		logger:  logger,
	}
}

func (p *Ping) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "plain/text")

	pingContext, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	if err := p.storage.Ping(pingContext); err != nil {
		p.logger.Warn("Ping database failed", zap.Error(err))
		http.Error(resp, "Lost connection Database", http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusOK)
}
