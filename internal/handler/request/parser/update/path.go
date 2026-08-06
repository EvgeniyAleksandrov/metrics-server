package update

import (
	"net/http"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/params"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/interfaces"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/types"
	models "github.com/EvgeniyAleksandrov/metrics-server/internal/model"
	"go.uber.org/zap"
)

type Path struct {
	logger interfaces.Logger
}

func NewPath(logger interfaces.Logger) *Path {
	return &Path{
		logger: logger,
	}
}

func (p *Path) Parse(req *http.Request) (*params.Update, error) {
	id := req.PathValue("name")
	mType := req.PathValue("method")
	value := req.PathValue("value")

	if id == "" || mType == "" || value == "" {
		return nil, handler.ErrInvalidParams
	}

	updateParams := &params.Update{
		ID:    id,
		MType: mType,
	}

	switch mType {
	case models.Gauge:
		gValue, err := types.GaugeFromString(value)
		if err != nil {
			p.logger.Warn("Parse gauge", zap.Error(err))
			return nil, handler.ErrInvalidValueFormat
		}

		gParamsValue := float64(gValue)

		updateParams.Value = &gParamsValue
	case models.Counter:
		cValue, err := types.CounterFromString(value)
		if err != nil {
			p.logger.Warn("Parse counter", zap.Error(err))
			return nil, handler.ErrInvalidValueFormat
		}

		cParamsValue := int64(cValue)

		updateParams.Delta = &cParamsValue
	default:
		return nil, handler.ErrUnsupportedMetricType
	}

	return updateParams, nil
}
