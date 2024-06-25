package usecase

import (
	"multi-site-dashboard-go/internal/domain"
	repo "multi-site-dashboard-go/internal/repository"
)

// In theory, the repository should return domain models.
// However, the persistence model can be different from the domain model.
// Hence, treat the persistence model as a DTO, and map it to the domain model
// in Use Cases.
type persistenceModelAdapter struct {}

var pma = persistenceModelAdapter{}

func (pma persistenceModelAdapter) MachineResourceUsage(arg repo.MachineResourceUsage) domain.MachineResourceUsage {
	return domain.MachineResourceUsage{
		Machine: arg.Machine,
		Metric1: arg.Metric1,
		Metric2: arg.Metric2,
		Metric3: arg.Metric3,
		CreatedAt: arg.CreatedAt.Time.String(),
	}
}

func (pma persistenceModelAdapter) AggMachineResourceUsage(arg repo.GetAggregatedMachineResourceUsageRow) domain.AggMachineResourceUsage {
	return domain.AggMachineResourceUsage{
		Metric1: arg.Metric1,
		Metric2: arg.Metric2,
		Metric3: arg.Metric3,
		Bucket: arg.Bucket.Time.String(),
	}
}
