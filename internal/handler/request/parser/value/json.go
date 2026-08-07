package value

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/params"
	models "github.com/EvgeniyAleksandrov/metrics-server/internal/model"
)

type JSONParams struct {
	maxBytes int64
}

func NewJSON(maxBytes int64) *JSONParams {
	if maxBytes <= 0 {
		panic("MaxBytes must be more than 0")
	}

	return &JSONParams{
		maxBytes: maxBytes,
	}
}

func (p *JSONParams) Parse(r *http.Request) (*params.Value, error) {
	if r == nil || r.Body == nil {
		return nil, handler.ErrInvalidParams
	}

	defer r.Body.Close()

	decoder := json.NewDecoder(
		http.MaxBytesReader(nil, r.Body, p.maxBytes),
	)

	metrics := &models.Metrics{}

	if err := decoder.Decode(metrics); err != nil {
		return nil, fmt.Errorf("decode json: %w", err)
	}

	return &params.Value{
		ID:    metrics.ID,
		MType: metrics.MType,
	}, nil
}
