package update

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/params"
	models "github.com/EvgeniyAleksandrov/metrics-server/internal/model"
)

var ErrInvalidRequest = errors.New("invalid request")

type JSON struct {
	maxBytes int64
}

func NewJSON(maxBytes int64) *JSON {
	if maxBytes <= 0 {
		panic("MaxBytes must be more than 0")
	}

	return &JSON{
		maxBytes: maxBytes,
	}
}

func (p *JSON) Parse(r *http.Request) ([]params.Update, error) {
	if r == nil || r.Body == nil {
		return nil, ErrInvalidRequest
	}

	defer r.Body.Close()

	decoder := json.NewDecoder(
		http.MaxBytesReader(nil, r.Body, p.maxBytes),
	)

	updateParamsList := make([]params.Update, 0)

	if err := decoder.Decode(&updateParamsList); err != nil {
		log.Println("PARSING ERROR:", err)
		return nil, fmt.Errorf("decode batch json: %w", err)
	}

	log.Printf("DATA: %+v", updateParamsList)

	for _, updateParams := range updateParamsList {
		if updateParams.ID == "" {
			return nil, handler.ErrInvalidParams
		}

		switch updateParams.MType {
		case models.Gauge:
			if updateParams.Value == nil {
				return nil, handler.ErrInvalidParams

			}
		case models.Counter:
			if updateParams.Delta == nil {
				return nil, handler.ErrInvalidParams
			}
		default:
			return nil, handler.ErrUnsupportedMetricType
		}
	}

	return updateParamsList, nil
}
