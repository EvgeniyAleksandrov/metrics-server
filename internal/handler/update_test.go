package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler"
	parserupdate "github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/parser/update"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/response/update"
	models "github.com/EvgeniyAleksandrov/metrics-server/internal/model"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/repository"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/service"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/service/method"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

const testMaxBytes = 8 * 1024

type SuitUpdate struct {
	suite.Suite
}

func TestSuitUpdate(t *testing.T) {
	t.Parallel()

	suite.Run(t, &SuitUpdate{})
}

func (s *SuitUpdate) TestUpdatePathValue_CorrectGaugeRequest_WriteOKStatus() {
	logger := zap.NewNop()
	memStorage := repository.NewMemStorage()

	processor := service.NewProcessor(method.NewGauge(memStorage, logger), method.NewCounter(memStorage, logger))

	req := s.buildRequestWithPathParams("Param1", "gauge", "12.22")
	rec := httptest.NewRecorder()

	sut := handler.NewUpdate(processor, parserupdate.NewPath(logger), update.NewValueWriter())

	sut.ServeHTTP(rec, req)

	s.Require().Equal(rec.Code, http.StatusOK)
	s.Require().Equal(rec.Header().Get("Content-Type"), "text/plain; charset=utf-8")

	value, err := memStorage.GetGauge("Param1")
	require.NoError(s.T(), err)
	require.Equal(s.T(), value, 12.22)
}

func (s *SuitUpdate) TestUpdatePathValue_CorrectCounterRequest_WriteOKStatus() {
	logger := zap.NewNop()
	memStorage := repository.NewMemStorage()

	processor := service.NewProcessor(method.NewGauge(memStorage, logger), method.NewCounter(memStorage, logger))

	req := s.buildRequestWithPathParams("Param1", "counter", "25")
	rec := httptest.NewRecorder()

	sut := handler.NewUpdate(processor, parserupdate.NewPath(logger), update.NewValueWriter())

	sut.ServeHTTP(rec, req)

	s.Require().Equal(rec.Code, http.StatusOK)
	s.Require().Equal(rec.Header().Get("Content-Type"), "text/plain; charset=utf-8")

	value, err := memStorage.GetCounter("Param1")
	require.NoError(s.T(), err)
	require.Equal(s.T(), value, int64(25))
}

func (s *SuitUpdate) TestUpdatePathValue_RequestWithUnsupportedMethod_WriteErrorStatus() {
	logger := zap.NewNop()
	memStorage := repository.NewMemStorage()

	processor := service.NewProcessor(method.NewGauge(memStorage, logger), method.NewCounter(memStorage, logger))

	req := s.buildRequestWithPathParams("Param1", "testMethod", "25")
	rec := httptest.NewRecorder()

	sut := handler.NewUpdate(processor, parserupdate.NewPath(logger), update.NewValueWriter())

	sut.ServeHTTP(rec, req)

	s.Require().Equal(rec.Code, http.StatusBadRequest)
	s.Require().Equal(rec.Header().Get("Content-Type"), "text/plain; charset=utf-8")

	gaugeValues, err := memStorage.GetAllGaugeValues()
	s.Require().NoError(err)
	s.Require().Empty(gaugeValues)

	counterValues, err := memStorage.GetAllGaugeValues()
	s.Require().NoError(err)
	s.Require().Empty(counterValues)
}

func (s *SuitUpdate) TestUpdatePathValue_InvalidValueFormat_WriteErrorStatus() {
	logger := zap.NewNop()
	memStorage := repository.NewMemStorage()

	processor := service.NewProcessor(method.NewGauge(memStorage, logger), method.NewCounter(memStorage, logger))

	req := s.buildRequestWithPathParams("Param1", "counter", "g")
	rec := httptest.NewRecorder()

	sut := handler.NewUpdate(processor, parserupdate.NewPath(logger), update.NewValueWriter())

	sut.ServeHTTP(rec, req)

	s.Require().Equal(rec.Code, http.StatusBadRequest)
	s.Require().Equal(rec.Header().Get("Content-Type"), "text/plain; charset=utf-8")

	gaugeValues, err := memStorage.GetAllGaugeValues()
	s.Require().NoError(err)
	s.Require().Empty(gaugeValues)

	counterValues, err := memStorage.GetAllGaugeValues()
	s.Require().NoError(err)
	s.Require().Empty(counterValues)
}

func (s *SuitUpdate) TestUpdatePathValue_RequestWithoutValue_WriteNotFoundErrorStatus() {
	logger := zap.NewNop()
	memStorage := repository.NewMemStorage()

	processor := service.NewProcessor(method.NewGauge(memStorage, logger), method.NewCounter(memStorage, logger))

	req := s.buildRequestWithOutValue("Param1", "counter")
	rec := httptest.NewRecorder()

	sut := handler.NewUpdate(processor, parserupdate.NewPath(logger), update.NewValueWriter())

	sut.ServeHTTP(rec, req)

	s.Require().Equal(rec.Code, http.StatusNotFound)
	s.Require().Equal(rec.Header().Get("Content-Type"), "text/plain; charset=utf-8")

	gaugeValues, err := memStorage.GetAllGaugeValues()
	s.Require().NoError(err)
	s.Require().Empty(gaugeValues)

	counterValues, err := memStorage.GetAllGaugeValues()
	s.Require().NoError(err)
	s.Require().Empty(counterValues)
}

func (s *SuitUpdate) TestUpdateJSONValue_CorrectGaugeRequest_WriteOKStatus() {
	logger := zap.NewNop()
	memStorage := repository.NewMemStorage()

	processor := service.NewProcessor(method.NewGauge(memStorage, logger), method.NewCounter(memStorage, logger))

	gaugeValue := float64(14.44)
	req := s.buildRequestWithJson(models.Metrics{
		ID:    "ParamJSON1",
		MType: "gauge",
		Value: &gaugeValue,
	})
	rec := httptest.NewRecorder()

	sut := handler.NewUpdate(processor, parserupdate.NewJSON(testMaxBytes), update.NewJSONWriter())

	sut.ServeHTTP(rec, req)

	s.Require().Equal(rec.Code, http.StatusOK)
	s.Require().Equal(rec.Header().Get("Content-Type"), "application/json; charset=utf-8")

	value, err := memStorage.GetGauge("ParamJSON1")
	require.NoError(s.T(), err)
	require.Equal(s.T(), value, 14.44)
}

func (s *SuitUpdate) TestUpdateJSONValue_CorrectCounterRequest_WriteOKStatus() {
	logger := zap.NewNop()
	memStorage := repository.NewMemStorage()

	processor := service.NewProcessor(method.NewGauge(memStorage, logger), method.NewCounter(memStorage, logger))

	delta := int64(44)
	req := s.buildRequestWithJson(models.Metrics{
		ID:    "ParamJSON1",
		MType: "counter",
		Delta: &delta,
	})
	rec := httptest.NewRecorder()

	sut := handler.NewUpdate(processor, parserupdate.NewJSON(testMaxBytes), update.NewJSONWriter())

	sut.ServeHTTP(rec, req)

	s.Require().Equal(rec.Code, http.StatusOK)
	s.Require().Equal(rec.Header().Get("Content-Type"), "application/json; charset=utf-8")

	value, err := memStorage.GetCounter("ParamJSON1")
	require.NoError(s.T(), err)
	require.Equal(s.T(), value, int64(44))
}

func (s *SuitUpdate) TestUpdateJSONValue_InvalidRequest_WriteOKStatus() {
	logger := zap.NewNop()
	memStorage := repository.NewMemStorage()

	processor := service.NewProcessor(method.NewGauge(memStorage, logger), method.NewCounter(memStorage, logger))

	gaugeValue := 44.1
	req := s.buildRequestWithJson(models.Metrics{
		ID:    "ParamJSON1",
		MType: "counter",
		Value: &gaugeValue,
	})
	rec := httptest.NewRecorder()

	sut := handler.NewUpdate(processor, parserupdate.NewJSON(testMaxBytes), update.NewJSONWriter())

	sut.ServeHTTP(rec, req)

	s.Require().Equal(rec.Code, http.StatusNotFound)
	s.Require().Equal(rec.Header().Get("Content-Type"), "application/json; charset=utf-8")

	gaugeValues, err := memStorage.GetAllGaugeValues()
	s.Require().NoError(err)
	s.Require().Empty(gaugeValues)

	counterValues, err := memStorage.GetAllGaugeValues()
	s.Require().NoError(err)
	s.Require().Empty(counterValues)
}

// HELPERS

func (s *SuitUpdate) buildRequestWithPathParams(name, method, value string) *http.Request {
	req, err := http.NewRequest(http.MethodPost, "", nil)
	s.Require().NoError(err)

	req.SetPathValue("name", name)
	req.SetPathValue("method", method)
	req.SetPathValue("value", value)

	return req
}

func (s *SuitUpdate) buildRequestWithOutValue(name, method string) *http.Request {
	req, err := http.NewRequest(http.MethodPost, "", nil)
	s.Require().NoError(err)

	req.SetPathValue("name", name)
	req.SetPathValue("method", method)

	return req
}

func (s *SuitUpdate) buildRequestWithJson(metric models.Metrics) *http.Request {
	data, err := json.Marshal(metric)
	s.Require().NoError(err)

	req, err := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(data))
	s.Require().NoError(err)

	return req
}
