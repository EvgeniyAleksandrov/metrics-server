//go:generate mockgen -source=publisher.go -destination=publisher_mock_test.go -package=metric_test
package metric

import (
	"fmt"
	"net/http"
)

var ErrNotPublished = fmt.Errorf("metrics not published")

type HTTPClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type Publisher struct {
	httpClient HTTPClient
}

func NewPublisher(httpClient HTTPClient) *Publisher {
	return &Publisher{
		httpClient: httpClient,
	}
}

func (p *Publisher) Publish(server string, m *Metric) error {
	req, err := p.buildUpdateRequest(server, m)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	res, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("make request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("publishing err with status %s: %w", res.Status, ErrNotPublished)
	}

	return nil
}

func (p *Publisher) buildUpdateRequest(server string, m *Metric) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, p.buildURL(server, m), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "text/plain")

	return req, nil
}

func (p *Publisher) buildURL(server string, m *Metric) string {
	return fmt.Sprintf("http://%s/update/%s/%s/%s", server, string(m.Type()), m.Name(), m.Value())
}
