package value

import (
	"net/http"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/params"
	models "github.com/EvgeniyAleksandrov/metrics-server/internal/model"
)

type Path struct{}

func NewPath() *Path {
	return &Path{}
}

func (p *Path) Parse(req *http.Request) (*params.Value, error) {
	id := req.PathValue("name")
	mType := req.PathValue("method")

	if id == "" || mType == "" {
		return nil, handler.ErrInvalidParams
	}

	if mType != models.Gauge && mType != models.Counter {
		return nil, handler.ErrUnsupportedMetricType
	}

	return &params.Value{
		ID:    id,
		MType: mType,
	}, nil
}
