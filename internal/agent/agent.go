package agent

import (
	"context"
	"log"
	"time"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/resource"
)

type ResourceManager interface {
	Update() error
	Get(resourceType resource.Type) (resource.Resource, error)
}

type MetricGetter interface {
	Resource() resource.Type
	GetAll(resource.Resource) ([]*metric.Metric, error)
}

type MetricPublisher interface {
	Publish(server string, m *metric.Metric) error
}
type Agent struct {
	resourceManager ResourceManager
	metricGetters   []MetricGetter
	metricPublisher MetricPublisher

	serverURL string

	pullInterval   time.Duration
	reportInterval time.Duration
}

func NewAgent(
	resourceManager ResourceManager,
	metricGetters []MetricGetter,
	metricPublisher MetricPublisher,
	serverURL string,
	pullInterval time.Duration,
	reportInterval time.Duration,
) *Agent {
	for _, mGetter := range metricGetters {
		resType := mGetter.Resource()
		if _, err := resourceManager.Get(resType); err != nil {
			log.Fatalf("Resource types %s not found: %s", resType, err.Error())
		}

		log.Printf("resource added: %s", resType)
	}

	return &Agent{
		resourceManager: resourceManager,
		metricGetters:   metricGetters,
		metricPublisher: metricPublisher,
		serverURL:       serverURL,
		pullInterval:    pullInterval,
		reportInterval:  reportInterval,
	}
}

func (a *Agent) Run(ctx context.Context) {
	pullTicker := time.NewTicker(a.pullInterval)
	defer pullTicker.Stop()

	publishDelayTicker := time.NewTicker(a.reportInterval)
	defer publishDelayTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-pullTicker.C:
			if err := a.resourceManager.Update(); err != nil {
				log.Printf("Update metrics error: %s", err.Error())
			}
		case <-publishDelayTicker.C:
			allMetrics := make([]*metric.Metric, 0)

			for _, mGetter := range a.metricGetters {
				res, err := a.resourceManager.Get(mGetter.Resource())
				if err != nil {
					log.Printf("Get resource error: %s", err.Error())
					continue
				}

				metrics, err := mGetter.GetAll(res)
				if err != nil {
					log.Printf("Get metrics error: %s", err.Error())
					continue
				}

				allMetrics = append(allMetrics, metrics...)
			}

			for _, m := range allMetrics {
				if err := a.metricPublisher.Publish(a.serverURL, m); err != nil {
					log.Printf("Publish metrics %s error: %s", m.Name(), err.Error())
				}
			}
		}
	}
}
