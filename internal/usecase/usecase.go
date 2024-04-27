package uc

import (
	"context"
	"multi-site-dashboard-go/internal/domain"
	repo "multi-site-dashboard-go/internal/repository"
)

var pma = persistenceModelAdapter{}

type UseCaseService struct {
	Repository repo.ExtQuerier
}

func NewUseCaseService(repo repo.ExtQuerier) *UseCaseService {
	return &UseCaseService{Repository: repo}
}

func (s *UseCaseService) GetMachineResourceUsageRT(ctx context.Context, machine string) ([]domain.MachineResourceUsage, error) {
	closure := s.Repository.ExecWithPgTx(ctx, func(qtx repo.Querier) (interface{}, error) {
		rv, err := qtx.GetMachineResourceUsage(ctx, machine)
		if err != nil {
			return nil, err
		}
		return rv, err
	})
	i, err := closure()
	if err != nil {
		return nil, err
	}

	pm := i.([]repo.MachineResourceUsage)
	rv := make([]domain.MachineResourceUsage, 0, len(pm))
	for _, x := range pm {
		rv = append(rv, pma.MachineResourceUsage(x))
	}
	return rv, nil
}

func (s *UseCaseService) GetAggMachineResourceUsageRT(ctx context.Context, arg *domain.GetAggMachineResourceUsageParams) ([]domain.AggMachineResourceUsage, error) {
	closure := s.Repository.ExecWithPgTx(ctx, func(qtx repo.Querier) (interface{}, error) {
		pmArg := repo.GetAggregatedMachineResourceUsageParams{
			Machine: arg.Machine,
			TimeBucket: arg.TimeBucket,
			LookBackPeriod: arg.LookBackPeriod,
		}
		rv, err := qtx.GetAggregatedMachineResourceUsage(ctx, pmArg)
		if err != nil {
			return nil, err
		}
		return rv, err
	})
	i, err := closure()
	if err != nil {
		return nil, err
	}

	pm := i.([]repo.GetAggregatedMachineResourceUsageRow)
	rv := make([]domain.AggMachineResourceUsage, 0, len(pm))
	for _, x := range pm {
		rv = append(rv, pma.AggMachineResourceUsage(x))
	}
	return rv, nil
}

func (s *UseCaseService) CreateMachineResourceUsage(ctx context.Context, arg *domain.MachineResourceUsage) (domain.MachineResourceUsage, error) {
	pmArg := repo.CreateMachineResourceUsageParams{Machine: arg.Machine, Metric1: arg.Metric1, Metric2: arg.Metric2, Metric3: arg.Metric3}
	i, err := s.Repository.CreateMachineResourceUsage(ctx, pmArg)
	if err != nil {
		return domain.MachineResourceUsage{}, err
	}
	rv := pma.MachineResourceUsage(i)
	return rv, nil
}

func (s *UseCaseService) TestSuccessPgTransaction(ctx context.Context, arg *domain.MachineResourceUsage) (domain.MachineResourceUsage, error) {
	closure := s.Repository.ExecWithPgTx(ctx, func(qtx repo.Querier) (interface{}, error) {	
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

	pm := i.(repo.MachineResourceUsage)
	rv := pma.MachineResourceUsage(pm)
	return rv, nil
}

func (s *UseCaseService) TestFailedPgTransaction(ctx context.Context, arg *domain.MachineResourceUsage) (domain.MachineResourceUsage, error) {
	closure := s.Repository.ExecWithPgTx(ctx, func(qtx repo.Querier) (interface{}, error) {	
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

	pm:= i.(repo.MachineResourceUsage)
	rv := pma.MachineResourceUsage(pm)
	return rv, nil
}