package metric_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/translator"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/translator/translation"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var ErrTestPublishing = errors.New("test publication error")

func TestPublisherUpdate_PublicateGaugeMetric_ReturnNoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHTTPClient := NewMockHTTPClient(ctrl)

	testMetric := metric.NewMetric(metric.Alloc, types.Gauge(123.2))

	mockHTTPClient.EXPECT().Do(gomock.Any()).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       http.NoBody,
		},
		nil,
	).Times(1)

	metricTranslator := translator.NewMetric(
		translation.NewGauge(),
		translation.NewCounter(),
	)

	sut := metric.NewPublisher(mockHTTPClient, metricTranslator)

	err := sut.Publish("localhost", testMetric)
	require.NoError(t, err)
}

func TestPublisherUpdate_PublicateCounterMetric_ReturnNoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHTTPClient := NewMockHTTPClient(ctrl)

	testMetric := metric.NewMetric(metric.PollCount, types.Counter(22))

	mockHTTPClient.EXPECT().Do(gomock.Any()).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       http.NoBody,
		},
		nil,
	).Times(1)

	metricTranslator := translator.NewMetric(
		translation.NewGauge(),
		translation.NewCounter(),
	)

	sut := metric.NewPublisher(mockHTTPClient, metricTranslator)

	err := sut.Publish("localhost", testMetric)
	require.NoError(t, err)
}

func TestPublisherUpdate_PublicateOnUnavailableServer_ReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHTTPClient := NewMockHTTPClient(ctrl)

	testMetric := metric.NewMetric(metric.PollCount, types.Counter(22))

	mockHTTPClient.EXPECT().Do(gomock.Any()).Return(nil, ErrTestPublishing).Times(1)

	metricTranslator := translator.NewMetric(
		translation.NewGauge(),
		translation.NewCounter(),
	)

	sut := metric.NewPublisher(mockHTTPClient, metricTranslator)

	err := sut.Publish("localhost", testMetric)
	require.ErrorIs(t, err, ErrTestPublishing)
}

func TestPublisherUpdate_PublicateWithNotOKStatus_ReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHTTPClient := NewMockHTTPClient(ctrl)

	testMetric := metric.NewMetric(metric.PollCount, types.Counter(22))

	mockHTTPClient.EXPECT().Do(gomock.Any()).Return(
		&http.Response{
			StatusCode: http.StatusMethodNotAllowed,
			Header:     make(http.Header),
			Body:       http.NoBody,
		},
		nil,
	).Times(1)

	metricTranslator := translator.NewMetric(
		translation.NewGauge(),
		translation.NewCounter(),
	)

	sut := metric.NewPublisher(mockHTTPClient, metricTranslator)

	err := sut.Publish("localhost", testMetric)
	require.ErrorIs(t, err, metric.ErrNotPublished)
	require.ErrorContains(t, err, "publishing err with status")
}
