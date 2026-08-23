package agent

import (
	"context"
	"sync"
	"time"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/resource"
	"go.uber.org/zap"
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
	Publish(ctx context.Context, server string, m *metric.Metric) error
}

type Logger interface {
	Fatal(msg string, fields ...zap.Field)
	Info(msg string, flags ...zap.Field)
	Warn(msg string, flags ...zap.Field)
}

type Hasher interface {
	MakeHash(msg string) string
}

type Agent struct {
	resourceManager ResourceManager
	metricGetters   []MetricGetter
	metricPublisher MetricPublisher

	serverURL string

	pullInterval   time.Duration
	reportInterval time.Duration
	numOfWorkers   int

	logger Logger
}

func NewAgent(
	resourceManager ResourceManager,
	metricGetters []MetricGetter,
	metricPublisher MetricPublisher,
	serverURL string,
	pullInterval time.Duration,
	reportInterval time.Duration,
	numOfWorkers int,
	logger Logger,
) *Agent {
	for _, mGetter := range metricGetters {
		resType := mGetter.Resource()
		if _, err := resourceManager.Get(resType); err != nil {
			logger.Fatal(
				"Resource types not found",
				zap.String("resourceType", string(resType)),
				zap.Error(err),
			)
		}

		logger.Info("Resource added", zap.String("resourceType", string(resType)))
	}

	return &Agent{
		resourceManager: resourceManager,
		metricGetters:   metricGetters,
		metricPublisher: metricPublisher,
		serverURL:       serverURL,
		pullInterval:    pullInterval,
		reportInterval:  reportInterval,
		numOfWorkers:    numOfWorkers,
		logger:          logger,
	}
}

func (a *Agent) Run(ctx context.Context) {
	metricsQueue := make(chan metric.Metric)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		a.runCollectorJob(ctx, metricsQueue)
	}()

	for i := 0; i < a.numOfWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			a.runPublisherJob(ctx, metricsQueue, i+1)
		}()
	}

	wg.Wait()

	a.logger.Info("Agent finished")
}

func (a *Agent) runPublisherJob(ctx context.Context, metricsQueue <-chan metric.Metric, workerNum int) {
	for {
		select {
		case <-ctx.Done():
			return
		case m := <-metricsQueue:
			if err := a.metricPublisher.Publish(ctx, a.serverURL, &m); err != nil {
				a.logger.Warn(
					"Publish metrics error",
					zap.String("metricName", string(m.Name())),
					zap.Int("numOfWorker", workerNum),
					zap.Error(err),
				)
			}
		}
	}
}

func (a *Agent) runCollectorJob(ctx context.Context, metricsQueue chan<- metric.Metric) {
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
				a.logger.Warn("Update metrics error", zap.Error(err))
			}
		case <-publishDelayTicker.C:
			allMetrics := make([]*metric.Metric, 0)

			for _, mGetter := range a.metricGetters {
				res, err := a.resourceManager.Get(mGetter.Resource())
				if err != nil {
					a.logger.Warn("Get resource error", zap.Error(err))
					continue
				}

				metrics, err := mGetter.GetAll(res)
				if err != nil {
					a.logger.Warn("Get metric error", zap.Error(err))
					continue
				}

				allMetrics = append(allMetrics, metrics...)
			}

			for _, m := range allMetrics {
				select {
				case <-ctx.Done():
					return
				default:
					metricsQueue <- *m

				}
			}
		}
	}
}
