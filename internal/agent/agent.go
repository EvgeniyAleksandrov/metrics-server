package agent

import (
	"fmt"
	"log"
	"time"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/resource"
)

var (
	updateDelay  = time.Second * 2
	publishDelay = time.Second * 10
	waitDelay    = time.Second * 1
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
	stop            bool

	serverURL string
}

func NewAgent(
	resourceManager ResourceManager,
	metricGetters []MetricGetter,
	metricPublisher MetricPublisher,
	serverURL string,
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
	}
}

func (a *Agent) Run() error {
	log.Println("Start agent.")

	updateTicker := time.NewTicker(updateDelay)
	defer updateTicker.Stop()

	publishDelayTicker := time.NewTicker(publishDelay)
	defer publishDelayTicker.Stop()

	for !a.stop {
		select {
		case <-updateTicker.C:
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
					return fmt.Errorf("get metrics %s: %w", mGetter.Resource(), err)
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

	log.Println("Stop agent.")

	return nil
}
