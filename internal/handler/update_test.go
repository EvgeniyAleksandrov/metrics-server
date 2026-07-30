package handler_test

import (
	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/params"
	models "github.com/EvgeniyAleksandrov/metrics-server/internal/model"
)

//import (
//	"errors"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//
//	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler"
//	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/params"
//	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/parser/update"
//	responseupdate "github.com/EvgeniyAleksandrov/metrics-server/internal/handler/response/update"
//	models "github.com/EvgeniyAleksandrov/metrics-server/internal/model"
//	"github.com/EvgeniyAleksandrov/metrics-server/internal/service"
//	"github.com/stretchr/testify/require"
//	"go.uber.org/mock/gomock"
//	"go.uber.org/zap"
//)

//func TestUpdate_ServeHTTPWithValueParser(t *testing.T) {
//	tests := map[string]struct {
//		httpMethod     string
//		pathMethod     string
//		pathName       string
//		pathValue      string
//		updateErr      error
//		expectedStatus int
//	}{
//		"success": {
//			httpMethod:     http.MethodPost,
//			pathMethod:     "gauge",
//			pathName:       "Alloc",
//			pathValue:      "10.5",
//			expectedStatus: http.StatusOK,
//		},
//		"unsupported metric method": {
//			httpMethod:     http.MethodPost,
//			pathMethod:     "unknown",
//			pathName:       "Alloc",
//			pathValue:      "10",
//			updateErr:      service.ErrUnsupportedProcessMethod,
//			expectedStatus: http.StatusBadRequest,
//		},
//		"invalid value": {
//			httpMethod:     http.MethodPost,
//			pathMethod:     "gauge",
//			pathName:       "Alloc",
//			pathValue:      "abc",
//			updateErr:      service.ErrInvalidValueFormat,
//			expectedStatus: http.StatusBadRequest,
//		},
//		"internal error": {
//			httpMethod:     http.MethodPost,
//			pathMethod:     "gauge",
//			pathName:       "Alloc",
//			pathValue:      "10",
//			updateErr:      errors.New("boom"),
//			expectedStatus: http.StatusInternalServerError,
//		},
//		"unexpected http method": {
//			httpMethod:     http.MethodGet,
//			expectedStatus: http.StatusMethodNotAllowed,
//		},
//	}
//
//	for name, tt := range tests {
//		t.Run(name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			mockProcessor := NewMockUpdateMetricProcessor(ctrl)
//			mockProcessor.EXPECT().Update(tt.pathMethod, tt.pathName, tt.pathValue).Return(tt.updateErr).AnyTimes()
//
//			req := httptest.NewRequest(tt.httpMethod, "/", nil)
//
//			req.SetPathValue("method", tt.pathMethod)
//			req.SetPathValue("name", tt.pathName)
//			req.SetPathValue("value", tt.pathValue)
//
//			rec := httptest.NewRecorder()
//
//			sut := handler.NewUpdate(mockProcessor, update.NewPath(zap.NewNop()), responseupdate.NewValueWriter())
//
//			sut.ServeHTTP(rec, req)
//
//			require.Equal(t, rec.Code, tt.expectedStatus)
//		})
//	}
//}

//func TestUpdate_ServeHTTPWithJsonParser(t *testing.T) {
//	tests := map[string]struct {
//		httpMethod     string
//		pathMethod     string
//		pathName       string
//		pathValue      string
//		jsonBody       string
//		updateErr      error
//		expectedStatus int
//	}{
//		"success": {
//			httpMethod:     http.MethodPost,
//			pathMethod:     "gauge",
//			pathName:       "Alloc",
//			pathValue:      "10.5",
//			jsonBody:       `{"id": "Alloc", "type":"gauge", "value":10.5}`,
//			expectedStatus: http.StatusOK,
//		},
//		//"unsupported metric method": {
//		//	httpMethod:     http.MethodPost,
//		//	pathMethod:     "unknown",
//		//	pathName:       "Alloc",
//		//	pathValue:      "10",
//		//	updateErr:      service.ErrUnsupportedProcessMethod,
//		//	expectedStatus: http.StatusBadRequest,
//		//},
//		//"invalid value": {
//		//	httpMethod:     http.MethodPost,
//		//	pathMethod:     "gauge",
//		//	pathName:       "Alloc",
//		//	pathValue:      "abc",
//		//	updateErr:      service.ErrInvalidValueFormat,
//		//	expectedStatus: http.StatusBadRequest,
//		//},
//		//"internal error": {
//		//	httpMethod:     http.MethodPost,
//		//	pathMethod:     "gauge",
//		//	pathName:       "Alloc",
//		//	pathValue:      "10",
//		//	updateErr:      errors.New("boom"),
//		//	expectedStatus: http.StatusInternalServerError,
//		//},
//		//"unexpected http method": {
//		//	httpMethod:     http.MethodGet,
//		//	expectedStatus: http.StatusMethodNotAllowed,
//		//},
//	}
//
//	for name, tt := range tests {
//		t.Run(name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			mockProcessor := NewMockUpdateMetricProcessor(ctrl)
//			mockProcessor.EXPECT().Update(tt.pathMethod, tt.pathName, tt.pathValue).Return(tt.updateErr).AnyTimes()
//
//			req := httptest.NewRequest(tt.httpMethod, "/", []byte(tt.jsonBody))
//
//			rec := httptest.NewRecorder()
//
//			sut := handler.NewUpdate(mockProcessor, update.NewPath())
//
//			sut.ServeHTTP(rec, req)
//
//			require.Equal(t, rec.Code, tt.expectedStatus)
//		})
//	}
//}

func buildGaugeUpdateRequest(name string, value float64) *params.Update {
	return &params.Update{
		ID:    name,
		MType: models.Gauge,
		Value: &value,
	}
}
