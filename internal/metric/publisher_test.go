package metric_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var ErrTestPublishing = errors.New("test publication error")

func TestPublisherUpdate_PublicateGaugeMetric_ReturnNoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHTTPClient := NewMockHTTPClient(ctrl)

	testMetric := metric.NewMetric(metric.Alloc, types.Gauge(123.2))

	expectedRequest := buildPublisherRequest("http://localhost/update/gauge/Alloc/123.2")

	mockHTTPClient.EXPECT().Do(expectedRequest).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       http.NoBody,
		},
		nil,
	).Times(1)

	sut := metric.NewPublisher(mockHTTPClient)

	err := sut.Publish("localhost", testMetric)
	require.NoError(t, err)
}

func TestPublisherUpdate_PublicateCounterMetric_ReturnNoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHTTPClient := NewMockHTTPClient(ctrl)

	testMetric := metric.NewMetric(metric.PullCount, types.Counter(22))

	expectedRequest := buildPublisherRequest("http://localhost/update/counter/PullCount/22")

	mockHTTPClient.EXPECT().Do(expectedRequest).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       http.NoBody,
		},
		nil,
	).Times(1)

	sut := metric.NewPublisher(mockHTTPClient)

	err := sut.Publish("localhost", testMetric)
	require.NoError(t, err)
}

func TestPublisherUpdate_PublicateOnUnavailableServer_ReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHTTPClient := NewMockHTTPClient(ctrl)

	testMetric := metric.NewMetric(metric.PullCount, types.Counter(22))

	expectedRequest := buildPublisherRequest("http://localhost/update/counter/PullCount/22")

	mockHTTPClient.EXPECT().Do(expectedRequest).Return(nil, ErrTestPublishing).Times(1)

	sut := metric.NewPublisher(mockHTTPClient)

	err := sut.Publish("localhost", testMetric)
	require.ErrorIs(t, err, ErrTestPublishing)
}

func TestPublisherUpdate_PublicateWithNotOKStatus_ReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHTTPClient := NewMockHTTPClient(ctrl)

	testMetric := metric.NewMetric(metric.PullCount, types.Counter(22))

	expectedRequest := buildPublisherRequest("http://localhost/update/counter/PullCount/22")

	mockHTTPClient.EXPECT().Do(expectedRequest).Return(
		&http.Response{
			StatusCode: http.StatusMethodNotAllowed,
			Header:     make(http.Header),
			Body:       http.NoBody,
		},
		nil,
	).Times(1)

	sut := metric.NewPublisher(mockHTTPClient)

	err := sut.Publish("localhost", testMetric)
	require.ErrorIs(t, err, metric.ErrNotPublished)
	require.ErrorContains(t, err, "publishing err with status")
}

func buildPublisherRequest(url string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, url, nil)
	req.Header.Set("Content-Type", "text/plain")

	return req
}
