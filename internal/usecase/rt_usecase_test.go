package usecase

import (
	"context"
	"errors"
	domain "multi-site-dashboard-go/internal/domain"
	"multi-site-dashboard-go/internal/repository"
	ucmock "multi-site-dashboard-go/internal/usecase/mock"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCaseGetAggMachineResourceUsageRT(t *testing.T) {
	// Setup usecase.
	ctx := context.Background()
	ucm := ucmock.ProvideUseCaseMocks(t)
	uc := NewUseCaseService(ucm.Repository, ucm.EventPublisher, ucm.Broadcaster)
	
	arg := &domain.GetAggMachineResourceUsageParams{
		Machine: "testMachine",
		TimeBucket: "1 day",
		LookBackPeriod: "1 week",
	}

	t.Run("should succeed", func(t *testing.T) {
		pmv := []repository.GetAggregatedMachineResourceUsageRow{
			{
				Metric1: 105,
				Metric2: 50,
				Metric3: 80,
			},
			{
				Metric1: 55,
				Metric2: 9,
				Metric3: 18,
			},
		} 
		ucm.Repository.EXPECT().GetAggregatedMachineResourceUsage(mock.Anything, mock.Anything).Return(pmv, nil).Once()
		rv, err := uc.GetAggMachineResourceUsageRT(ctx, arg)
		assert.Nil(t, err)
		assert.Equal(t, len(rv), 2)
		assert.Equal(t, rv[0].Metric1, float64(105))
		assert.Equal(t, rv[1].Metric3, float64(18))
	})

	t.Run("should return error on failure", func(t *testing.T) {
		ucm.Repository.EXPECT().GetAggregatedMachineResourceUsage(mock.Anything, mock.Anything).Return(nil, ucmock.ErrDB).Once()
		rv, err := uc.GetAggMachineResourceUsageRT(ctx, arg)
		assert.Nil(t, rv)
		assert.True(t, errors.Is(err, ucmock.ErrDB))
	})
}

func TestUseCaseCreateMachineResourceUsageRT(t *testing.T) {
	// Setup usecase.
	ctx := context.Background()
	ucm := ucmock.ProvideUseCaseMocks(t)
	uc := NewUseCaseService(ucm.Repository, ucm.EventPublisher, ucm.Broadcaster)

	metric1 := int32(105)
	metric2 := int32(70)
	metric3 := int32(61)
	arg := &domain.CreateMachineResourceUsageParams{
		Machine: "testMachine",
		Metric1: &metric1,
		Metric2: &metric2,
		Metric3: &metric3,
	}
	pmv := repository.MachineResourceUsage{
		Machine: "testMachine",
		Metric1: metric1,
		Metric2: metric2,
		Metric3: metric3,
	} 
	t.Run("should succeed", func(t *testing.T) {
		ucm.Repository.EXPECT().CreateMachineResourceUsage(mock.Anything, mock.Anything).Return(pmv, nil).Once()
		ucm.EventPublisher.EXPECT().PublishMachineResourceUsageEvent(mock.Anything, mock.Anything).Return(nil).Once()
		rv, err := uc.CreateMachineResourceUsageRT(ctx, arg)
		assert.Nil(t, err)
		assert.Equal(t, rv.Machine, "testMachine")
		assert.Equal(t, rv.Metric3, int32(61))
	})

	t.Run("should return error if publish failure", func(t *testing.T) {
		ucm.Repository.EXPECT().CreateMachineResourceUsage(mock.Anything, mock.Anything).Return(pmv, nil).Once()
		ucm.EventPublisher.EXPECT().PublishMachineResourceUsageEvent(mock.Anything, mock.Anything).Return(ucmock.ErrEventPublisher).Once()
		_, err := uc.CreateMachineResourceUsageRT(ctx, arg)
		assert.True(t, errors.Is(err, ucmock.ErrEventPublisher))
	})

	t.Run("should return error if create failure", func(t *testing.T) {
		ucm.Repository.EXPECT().CreateMachineResourceUsage(mock.Anything, mock.Anything).Return(repository.MachineResourceUsage{}, ucmock.ErrDB).Once()
		_, err := uc.CreateMachineResourceUsageRT(ctx, arg)
		assert.True(t, errors.Is(err, ucmock.ErrDB))
	})
}

func TestUseCaseCreateMachineResourceUsageAndBroadcastRT(t *testing.T) {
	ctx := context.Background()
	ucm := ucmock.ProvideUseCaseMocks(t)
	uc := NewUseCaseService(ucm.Repository, ucm.EventPublisher, ucm.Broadcaster)

	metricSentinel := int32(105)
	pmv := repository.MachineResourceUsage{} 
	arg := &domain.CreateMachineResourceUsageParams{
		Machine: "testMachine",
		Metric1: &metricSentinel,
		Metric2: &metricSentinel,
		Metric3: &metricSentinel,
	}
	t.Run("should succeed", func(t *testing.T) {
		ucm.Repository.EXPECT().CreateMachineResourceUsage(mock.Anything, mock.Anything).Return(pmv, nil).Once()
		ucm.Broadcaster.EXPECT().Broadcast(mock.Anything, mock.Anything).Return(nil).Once()
		_, err := uc.CreateMachineResourceUsageAndBroadcastRT(ctx, arg)
		assert.Nil(t, err)
	})

	t.Run("should return error if publish failure", func(t *testing.T) {
		ucm.Repository.EXPECT().CreateMachineResourceUsage(mock.Anything, mock.Anything).Return(pmv, nil).Once()
		ucm.Broadcaster.EXPECT().Broadcast(mock.Anything, mock.Anything).Return(ucmock.ErrBroadcast).Once()
		_, err := uc.CreateMachineResourceUsageAndBroadcastRT(ctx, arg)
		assert.True(t, errors.Is(err, ucmock.ErrBroadcast))
	})
}