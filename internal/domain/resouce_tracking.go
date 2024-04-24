package domain

type MachineResource struct {
	Machine string `json:"name" validate:"required"`
	Metric1 int `json:"metric1" validate:"required"`
	Metric2 int `json:"metric2" validate:"required"`
	Metric3 int `json:"metric3" validate:"required"`
}
