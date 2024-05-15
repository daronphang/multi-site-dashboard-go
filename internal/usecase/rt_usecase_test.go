package usecase

import (
	"context"
	"errors"
	domain "multi-site-dashboard-go/internal/domain"
	"multi-site-dashboard-go/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAggMachineResourceUsageRT(t *testing.T) {
	ctx := context.Background()
	ucm := provideUseCaseMocks(t)
	
	arg := &domain.GetAggMachineResourceUsageParams{
		Machine: "testMachine",
		TimeBucket: "1 day",
		LookBackPeriod: "1 week",
	}

	t.Run("query success", func(t *testing.T) {
		pmv := []repository.GetAggregatedMachineResourceUsageRow{{},{}} 
		ucm.repository.EXPECT().GetAggregatedMachineResourceUsage(mock.Anything, mock.Anything).Return(pmv, nil).Once()
		rv, err := ucm.usecase.GetAggMachineResourceUsageRT(ctx, arg)
		assert.Nil(t, err)
		assert.Equal(t, len(rv), 2)
	})

	t.Run("query failure", func(t *testing.T) {
		ucm.repository.EXPECT().GetAggregatedMachineResourceUsage(mock.Anything, mock.Anything).Return(nil, errDB).Once()
		rv, err := ucm.usecase.GetAggMachineResourceUsageRT(ctx, arg)
		assert.Nil(t, rv)
		assert.True(t, errors.Is(err, errDB))
	})
}

func TestCreateMachineResourceUsageRT(t *testing.T) {
	ctx := context.Background()
	ucm := provideUseCaseMocks(t)
	metricSentinel := int32(105)
	pmv := repository.MachineResourceUsage{} 
	
	arg := &domain.CreateMachineResourceUsageParams{
		Machine: "testMachine",
		Metric1: &metricSentinel,
		Metric2: &metricSentinel,
		Metric3: &metricSentinel,
	}
	t.Run("create and publish success", func(t *testing.T) {
		ucm.repository.EXPECT().CreateMachineResourceUsage(mock.Anything, mock.Anything).Return(pmv, nil).Once()
		ucm.eventPublisher.EXPECT().PublishMachineResourceUsageEvent(mock.Anything, mock.Anything).Return(nil).Once()
		_, err := ucm.usecase.CreateMachineResourceUsageRT(ctx, arg)
		assert.Nil(t, err)
	})

	t.Run("create success but publish failure", func(t *testing.T) {
		ucm.repository.EXPECT().CreateMachineResourceUsage(mock.Anything, mock.Anything).Return(pmv, nil).Once()
		ucm.eventPublisher.EXPECT().PublishMachineResourceUsageEvent(mock.Anything, mock.Anything).Return(errEventPublisher).Once()
		_, err := ucm.usecase.CreateMachineResourceUsageRT(ctx, arg)
		assert.True(t, errors.Is(err, errEventPublisher))
	})

	t.Run("create failure", func(t *testing.T) {
		ucm.repository.EXPECT().CreateMachineResourceUsage(mock.Anything, mock.Anything).Return(repository.MachineResourceUsage{}, errDB).Once()
		_, err := ucm.usecase.CreateMachineResourceUsageRT(ctx, arg)
		assert.True(t, errors.Is(err, errDB))
	})
}

func TestCreateMachineResourceUsageAndBroadcastRT(t *testing.T) {
	ctx := context.Background()
	ucm := provideUseCaseMocks(t)
	metricSentinel := int32(105)
	pmv := repository.MachineResourceUsage{} 
	
	arg := &domain.CreateMachineResourceUsageParams{
		Machine: "testMachine",
		Metric1: &metricSentinel,
		Metric2: &metricSentinel,
		Metric3: &metricSentinel,
	}
	t.Run("create and broadcast success", func(t *testing.T) {
		ucm.repository.EXPECT().CreateMachineResourceUsage(mock.Anything, mock.Anything).Return(pmv, nil).Once()
		ucm.broadcaster.EXPECT().Broadcast(mock.Anything, mock.Anything).Return(nil).Once()
		_, err := ucm.usecase.CreateMachineResourceUsageAndBroadcastRT(ctx, arg)
		assert.Nil(t, err)
	})

	t.Run("create success but publish failure", func(t *testing.T) {
		ucm.repository.EXPECT().CreateMachineResourceUsage(mock.Anything, mock.Anything).Return(pmv, nil).Once()
		ucm.broadcaster.EXPECT().Broadcast(mock.Anything, mock.Anything).Return(errBroadcast).Once()
		_, err := ucm.usecase.CreateMachineResourceUsageAndBroadcastRT(ctx, arg)
		assert.True(t, errors.Is(err, errBroadcast))
	})
}