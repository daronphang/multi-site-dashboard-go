package domain

// swagger:parameters CreateMachineResourceUsage
type CreateMachineResourceUsageParams struct {
	// in: body
	Machine string `json:"machine" validate:"required"`
	Metric1 *int32  `json:"metric1" validate:"required"`
	Metric2 *int32  `json:"metric2" validate:"required"`
	Metric3 *int32  `json:"metric3" validate:"required"`
}

// swagger:model MachineResourceUsage
type MachineResourceUsage struct {
	Machine string `json:"machine"`
	Metric1 int32  `json:"metric1"`
	Metric2 int32  `json:"metric2"`
	Metric3 int32  `json:"metric3"`
	CreatedAt string `json:"createdAt"`
}

// swagger:parameters GetAggMachineResourceUsage
type GetAggMachineResourceUsageParams struct {
	// in: path
	Machine        string `json:"machine" validate:"required"`
	// in: query
	TimeBucket     string `json:"timeBucket" validate:"required"`
	// in: query
	LookBackPeriod string `json:"lookBackPeriod" validate:"required"`
}

// swagger:model AggMachineResourceUsage
type AggMachineResourceUsage struct {
	Bucket string `json:"bucket"`
	Metric1 float64  `json:"metric1"`
	Metric2 float64  `json:"metric2"`
	Metric3 float64  `json:"metric3"`
}

type UpdateMachineResourceUsageParams struct {
	Metric1 *int32  `json:"metric1" validate:"required"`
	Machine string `json:"machine" validate:"required"`
}