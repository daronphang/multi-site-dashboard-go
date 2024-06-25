package usecase

import (
	"context"
	"encoding/json"
	"multi-site-dashboard-go/internal/domain"
	repo "multi-site-dashboard-go/internal/repository"
)

func (uc *UseCaseService) GetAggMachineResourceUsageRT(ctx context.Context, arg *domain.GetAggMachineResourceUsageParams) ([]domain.AggMachineResourceUsage, error) {
	pmArg := repo.GetAggregatedMachineResourceUsageParams{
		Machine: arg.Machine,
		TimeBucket: arg.TimeBucket,
		LookBackPeriod: arg.LookBackPeriod,
	}
	pmv, err := uc.Repository.GetAggregatedMachineResourceUsage(ctx, pmArg)
	if err != nil {
		return nil, err
	}
	rv := make([]domain.AggMachineResourceUsage, 0, len(pmv))
	for _, x := range pmv {
		rv = append(rv, pma.AggMachineResourceUsage(x))
	}
	return rv, nil
}

func (uc *UseCaseService) CreateMachineResourceUsageRT(ctx context.Context, arg *domain.CreateMachineResourceUsageParams) (domain.MachineResourceUsage, error) {
	pmArg := repo.CreateMachineResourceUsageParams{Machine: arg.Machine, Metric1: *arg.Metric1, Metric2: *arg.Metric2, Metric3: *arg.Metric3}
	pmv, err := uc.Repository.CreateMachineResourceUsage(ctx, pmArg)
	if err != nil {
		return domain.MachineResourceUsage{}, err
	}
	rv := pma.MachineResourceUsage(pmv)

	// TODO: testing
	if err := uc.EventPublisher.PublishMachineResourceUsageEvent(ctx, *arg); err != nil {
		return domain.MachineResourceUsage{}, err
	}
	return rv, nil
}

func (uc *UseCaseService) CreateMachineResourceUsageAndBroadcastRT(ctx context.Context, arg *domain.CreateMachineResourceUsageParams) (domain.MachineResourceUsage, error) {
	pmArg := repo.CreateMachineResourceUsageParams{Machine: arg.Machine, Metric1: *arg.Metric1, Metric2: *arg.Metric2, Metric3: *arg.Metric3}
	pmv, err := uc.Repository.CreateMachineResourceUsage(ctx, pmArg)
	if err != nil {
		return domain.MachineResourceUsage{}, err
	}
	rv := pma.MachineResourceUsage(pmv)

	data, err := json.Marshal(arg)
	if err != nil {
		return domain.MachineResourceUsage{}, err
	}
	if err := uc.Broadcaster.Broadcast(ctx, data); err != nil {
		return domain.MachineResourceUsage{}, err
	}
	return rv, nil
}

func (uc *UseCaseService) TestSuccessPgTransaction(ctx context.Context, arg *domain.MachineResourceUsage) (domain.MachineResourceUsage, error) {
	closure := uc.Repository.ExecWithPgTx(ctx, func(qtx repo.Querier) (interface{}, error) {	
		pmArg := repo.CreateMachineResourceUsageParams{Machine: arg.Machine, Metric1: arg.Metric1, Metric2: arg.Metric2, Metric3: arg.Metric3}	
		rv, err := qtx.CreateMachineResourceUsage(ctx, pmArg)
		if err != nil {
			return nil, err
		}
		
		up := repo.UpdateMachineResourceUsageParams{Machine: arg.Machine, Metric1: 99}
		err = qtx.UpdateMachineResourceUsage(ctx, up)
		if err != nil {
			return nil, err
		}
		return rv, err
	})
	i, err := closure()
	if err != nil {
		return domain.MachineResourceUsage{}, err
	}

	pmv := i.(repo.MachineResourceUsage)
	rv := pma.MachineResourceUsage(pmv)
	return rv, nil
}

func (uc *UseCaseService) TestFailedPgTransaction(ctx context.Context, arg *domain.MachineResourceUsage) (domain.MachineResourceUsage, error) {
	closure := uc.Repository.ExecWithPgTx(ctx, func(qtx repo.Querier) (interface{}, error) {	
		pmArg := repo.CreateMachineResourceUsageParams{Machine: arg.Machine, Metric1: arg.Metric1, Metric2: arg.Metric2, Metric3: arg.Metric3}	
		rv, err := qtx.CreateMachineResourceUsage(ctx, pmArg)
		if err != nil {
			return nil, err
		}
		
		up := repo.UpdateMachineResourceUsageParams{Machine: arg.Machine, Metric1: -1}
		err = qtx.UpdateMachineResourceUsage(ctx, up)
		if err != nil {
			return nil, err
		}
		return rv, err
	})
	i, err := closure()
	if err != nil {
		return domain.MachineResourceUsage{}, err
	}

	pmv := i.(repo.MachineResourceUsage)
	rv := pma.MachineResourceUsage(pmv)
	return rv, nil
}