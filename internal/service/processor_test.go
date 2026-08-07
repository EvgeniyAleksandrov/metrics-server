package service_test

//import (
//	"errors"
//	"testing"
//
//	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
//	"github.com/EvgeniyAleksandrov/metrics-server/internal/service"
//	"github.com/stretchr/testify/require"
//	"go.uber.org/mock/gomock"
//)
//
//var ErrTest = errors.New("test error")
//
//func TestNewProcessor_CreateNewProcessorWithAnyDifferentMethods_NoPanic(t *testing.T) {
//	ctrl := gomock.NewController(t)
//
//	mockMethod1 := NewMockMethod(ctrl)
//	mockMethod1.EXPECT().GetType().Return(metric.Counter).Times(1)
//
//	mockMethod2 := NewMockMethod(ctrl)
//	mockMethod2.EXPECT().GetType().Return(metric.Gauge).Times(1)
//
//	require.NotPanics(t, func() {
//		_ = service.NewProcessor(mockMethod1, mockMethod2)
//	})
//}
//
//func TestNewProcessor_CreateNewProcessorWithAnyIdenticalMethods_Panic(t *testing.T) {
//	ctrl := gomock.NewController(t)
//
//	mockMethod1 := NewMockMethod(ctrl)
//	mockMethod1.EXPECT().GetType().Return(metric.Gauge).Times(1)
//
//	mockMethod2 := NewMockMethod(ctrl)
//	mockMethod2.EXPECT().GetType().Return(metric.Gauge).Times(1)
//
//	require.Panics(t, func() {
//		_ = service.NewProcessor(mockMethod1, mockMethod2)
//	})
//}
//
//func TestProcessor_UpdateSupportedMetricData_NoErrorUpdate(t *testing.T) {
//	ctrl := gomock.NewController(t)
//
//	mockMethod1 := NewMockMethod(ctrl)
//	mockMethod1.EXPECT().GetType().Return(metric.Gauge).Times(1)
//	mockMethod1.EXPECT().Update("A", "12.5").Return(nil).Times(1)
//
//	sut := service.NewProcessor(mockMethod1)
//
//	err := sut.Update("gauge", "A", "12.5")
//	require.NoError(t, err)
//}
//
//func TestProcessor_UpdateNotSupportedMetricData_NoErrorUpdate(t *testing.T) {
//	ctrl := gomock.NewController(t)
//
//	mockMethod1 := NewMockMethod(ctrl)
//	mockMethod1.EXPECT().GetType().Return(metric.Gauge).Times(1)
//
//	sut := service.NewProcessor(mockMethod1)
//
//	err := sut.Update("not_supported", "A", "12.5")
//	require.ErrorIs(t, err, service.ErrUnsupportedProcessMethod)
//}
//
//func TestProcessor_UpdateDataWithError_ErrorUpdate(t *testing.T) {
//	ctrl := gomock.NewController(t)
//
//	mockMethod1 := NewMockMethod(ctrl)
//	mockMethod1.EXPECT().GetType().Return(metric.Gauge).Times(1)
//	mockMethod1.EXPECT().Update("A", "12.5").Return(ErrTest).Times(1)
//
//	sut := service.NewProcessor(mockMethod1)
//
//	err := sut.Update("gauge", "A", "12.5")
//	require.ErrorIs(t, err, ErrTest)
//}
