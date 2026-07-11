package resource_test

import (
	"testing"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/resource"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestManager_GetAvailableResourceWithOneResourceInList_ReturnTargetResource(t *testing.T) {
	sut := resource.NewManager(resource.NewMemory())

	res, err := sut.Get(resource.MemoryType)
	require.NoError(t, err)
	require.Equal(t, res.Type(), resource.MemoryType)
}

func TestManager_GetAvailableResourceWithAnyResourceInList_ReturnTargetResource(t *testing.T) {
	sut := resource.NewManager(resource.NewMemory(), resource.NewRandom())

	res, err := sut.Get(resource.RandomType)
	require.NoError(t, err)
	require.Equal(t, res.Type(), resource.RandomType)
}

func TestManager_GetUnavailableResource_ResourceNotFoundError(t *testing.T) {
	sut := resource.NewManager(resource.NewMemory(), resource.NewRandom())

	res, err := sut.Get(resource.PullCounterType)
	require.ErrorIs(t, err, resource.ErrResourceNotFound)
	require.Nil(t, res)
}

func TestManager_UpdateOneResource_ResourceWasUpdated(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockResource := NewMockResource(ctrl)
	mockResource.EXPECT().Type().Return(resource.Type("testType")).Times(2)
	mockResource.EXPECT().Update().Return(nil).Times(1)

	sut := resource.NewManager(mockResource)

	err := sut.Update()
	require.NoError(t, err)
}

func TestManager_UpdateAnyResource_ResourceWasUpdated(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockResource1 := NewMockResource(ctrl)
	mockResource1.EXPECT().Type().Return(resource.Type("testType1")).Times(2)
	mockResource1.EXPECT().Update().Return(nil).Times(1)

	mockResource2 := NewMockResource(ctrl)
	mockResource2.EXPECT().Type().Return(resource.Type("testType2")).Times(2)
	mockResource2.EXPECT().Update().Return(nil).Times(1)

	sut := resource.NewManager(mockResource1, mockResource2)

	err := sut.Update()
	require.NoError(t, err)
}

func TestManager_CreateWithAnyIdenticalResources_CreatePanic(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockResource1 := NewMockResource(ctrl)
	mockResource1.EXPECT().Type().Return(resource.Type("testType")).Times(2)

	mockResource2 := NewMockResource(ctrl)
	mockResource2.EXPECT().Type().Return(resource.Type("testType")).Times(1)

	require.Panics(t, func() {
		_ = resource.NewManager(mockResource1, mockResource2)
	})
}
