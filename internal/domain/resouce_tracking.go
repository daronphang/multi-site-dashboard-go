package domain

type MachineResourceUsage struct {
	Machine string `json:"machine" validate:"required"`
	Metric1 int32  `json:"metric1" validate:"int_required"`
	Metric2 int32  `json:"metric2" validate:"int_required"`
	Metric3 int32  `json:"metric3" validate:"int_required"`
	CreatedAt string `json:"createdAt"`
}

type AggMachineResourceUsage struct {
	Bucket string `json:"bucket"`
	Metric1 float64  `json:"metric1"`
	Metric2 float64  `json:"metric2"`
	Metric3 float64  `json:"metric3"`
}

type GetAggMachineResourceUsageParams struct {
	Machine        string `json:"machine" validate:"required"`
	TimeBucket     string `json:"timeBucket" validate:"required"`
	LookBackPeriod string `json:"lookBackPeriod" validate:"required"`
}

type UpdateMachineResourceUsageParams struct {
	Metric1 int32  `json:"metric1" validate:"required"`
	Machine string `json:"machine" validate:"required"`
}

