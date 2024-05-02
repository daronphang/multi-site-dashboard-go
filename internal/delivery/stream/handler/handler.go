package handler

import uc "multi-site-dashboard-go/internal/usecase"


type Handler struct {
	UseCase *uc.UseCaseService
}

func NewHandler(uc *uc.UseCaseService) *Handler {
	return &Handler{UseCase: uc}
}