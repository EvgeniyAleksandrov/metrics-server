package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/service"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUpdate_ServeHTTP(t *testing.T) {
	tests := map[string]struct {
		httpMethod     string
		pathMethod     string
		pathName       string
		pathValue      string
		updateErr      error
		expectedStatus int
	}{
		"success": {
			httpMethod:     http.MethodPost,
			pathMethod:     "gauge",
			pathName:       "Alloc",
			pathValue:      "10.5",
			expectedStatus: http.StatusOK,
		},
		"unsupported metric method": {
			httpMethod:     http.MethodPost,
			pathMethod:     "unknown",
			pathName:       "Alloc",
			pathValue:      "10",
			updateErr:      service.ErrUnsupportedProcessMethod,
			expectedStatus: http.StatusBadRequest,
		},
		"invalid value": {
			httpMethod:     http.MethodPost,
			pathMethod:     "gauge",
			pathName:       "Alloc",
			pathValue:      "abc",
			updateErr:      service.ErrInvalidValueFormat,
			expectedStatus: http.StatusBadRequest,
		},
		"internal error": {
			httpMethod:     http.MethodPost,
			pathMethod:     "gauge",
			pathName:       "Alloc",
			pathValue:      "10",
			updateErr:      errors.New("boom"),
			expectedStatus: http.StatusInternalServerError,
		},
		"unexpected http method": {
			httpMethod:     http.MethodGet,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockProcessor := NewMockUpdateMetricProcessor(ctrl)
			mockProcessor.EXPECT().Update(tt.pathMethod, tt.pathName, tt.pathValue).Return(tt.updateErr).AnyTimes()

			req := httptest.NewRequest(tt.httpMethod, "/", nil)

			req.SetPathValue("method", tt.pathMethod)
			req.SetPathValue("name", tt.pathName)
			req.SetPathValue("value", tt.pathValue)

			rec := httptest.NewRecorder()

			sut := handler.NewUpdate(mockProcessor)

			sut.ServeHTTP(rec, req)

			require.Equal(t, rec.Code, tt.expectedStatus)
		})
	}
}
