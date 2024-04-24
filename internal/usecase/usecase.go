package uc

import (
	"context"
	"multi-site-dashboard-go/internal"
	"multi-site-dashboard-go/internal/repository"
)

type UseCaseService struct {}

func NewUseCaseService() *UseCaseService {
	return &UseCaseService{}
}

func (s *UseCaseService) GetTimeSeriesMachineResourceUsageRT(ctx context.Context, machine string) error {
	db, err := internal.WirePgConnPool(ctx)
	if err != nil {
		return err
	}
	queries := repository.New(db)

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := queries.WithTx(tx)
	r, err := qtx.GetMachineResourceUsage(ctx, machine)
	if err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}