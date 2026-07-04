//go:generate mockgen -source=update.go -destination=update_mock_test.go -package=handler_test
package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/service/update"
)

type MetricProcessor interface {
	Update(method string, name string, value string) error
}

type Update struct {
	metricProcessor MetricProcessor
}

func NewUpdate(metricProcessor MetricProcessor) *Update {
	return &Update{
		metricProcessor: metricProcessor,
	}
}

func (u *Update) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(resp, "unsupported method", http.StatusMethodNotAllowed)
	}

	method := req.PathValue("method")
	name := req.PathValue("name")
	value := req.PathValue("value")

	if method == "" || name == "" || value == "" {
		http.Error(resp, "Not Found", http.StatusNotFound)
	}

	if err := u.metricProcessor.Update(method, name, value); err != nil {
		switch {
		case errors.Is(err, update.ErrUnsupportedProcessMethod):
			http.Error(resp, "Unsupported method", http.StatusBadRequest)

		case errors.Is(err, update.ErrInvalidValueFormat):
			http.Error(resp, "Invalid value format", http.StatusBadRequest)

		default:
			http.Error(resp, "Not initialized error", http.StatusInternalServerError)
		}

		return
	}

	log.Printf("metirc name:%s, value: %s, method: %s processed\n", name, value, method)
	u.writeHeaders(resp)
	resp.WriteHeader(http.StatusOK)

}

func (u *Update) writeHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "plain/text")
	w.Header().Set("charset", "utf-8")
}

// func (c *Update) POST(w http.ResponseWriter, r *http.Request) {
// 	param, err := c.paramsGetter.Get(r)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, metric.ErrEmptyName) || errors.Is(err, metric.ErrEmptyValue):
// 			log.Printf("counter POST: %s", err.Error())
// 			http.Error(w, "empty name", http.StatusNotFound)
// 		default:
// 			http.Error(w, "unexpected error", http.StatusBadRequest)
// 		}
// 	}

// 	i, err := strconv.Atoi(param.Value)
// 	if err != nil {
// 		http.Error(w, "not correct value format", http.StatusBadRequest)
// 		return
// 	}

// 	value := int64(i)

// 	// TODO: write param to data base

// 	w.Write([]byte(fmt.Sprintf("name: %s, value: %d", param.Name, value)))
// }

// func (u *Update) parseParams(r *http.Request) (*updateParams, error) {
// }
