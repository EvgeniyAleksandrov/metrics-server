//go:generate mockgen -source=publisher.go -destination=publisher_mock_test.go -package=metric_test
package metric

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"

	models "github.com/EvgeniyAleksandrov/metrics-server/internal/model"
)

var ErrNotPublished = fmt.Errorf("metrics not published")

type HTTPClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type Translator interface {
	Translate(m *Metric) (*models.Metrics, error)
}

type Publisher struct {
	httpClient HTTPClient
	translator Translator
}

func NewPublisher(httpClient HTTPClient, translator Translator) *Publisher {
	return &Publisher{
		httpClient: httpClient,
		translator: translator,
	}
}

func (p *Publisher) Publish(server string, m *Metric) error {
	req, err := p.buildUpdateRequest(server, m)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	_ = req.Body.Close()

	res, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("make request: %w", err)
	}

	if res == nil || res.Body == nil {
		return fmt.Errorf("empty response")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("publishing err with status %s: %w", res.Status, ErrNotPublished)
	}

	return nil
}

func (p *Publisher) buildUpdateRequest(server string, m *Metric) (*http.Request, error) {
	jsonData, err := p.buildJSON(m)
	if err != nil {
		return nil, fmt.Errorf("build json data: %w", err)
	}

	var byteBuffer bytes.Buffer

	gzWriter := gzip.NewWriter(&byteBuffer)

	if _, err := gzWriter.Write(jsonData); err != nil {
		return nil, fmt.Errorf("gzip write data: %w", err)
	}

	if err := gzWriter.Close(); err != nil {
		return nil, fmt.Errorf("close gzip writer: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://%s/update/", server),
		bytes.NewReader(byteBuffer.Bytes()),
	)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "application/gzip")
	req.Header.Set("Accept-Encoding", "gzip")

	return req, nil
}

func (p *Publisher) buildJSON(m *Metric) ([]byte, error) {
	metrics, err := p.translator.Translate(m)
	if err != nil {
		return nil, fmt.Errorf("translate metric: %w", err)
	}

	jsonData, err := json.Marshal(metrics)
	if err != nil {
		return nil, fmt.Errorf("marshaling metrics: %w", err)
	}

	return jsonData, nil
}
