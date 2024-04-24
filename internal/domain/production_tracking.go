package domain

type PTSalesOrder struct {
	Name string `json:"name" validate:"required"`
	ID int `json:"id" validate:"required"`
}