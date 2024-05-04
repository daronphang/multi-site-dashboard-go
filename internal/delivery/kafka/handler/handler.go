package handler

import uc "multi-site-dashboard-go/internal/usecase"


type KafkaHandler struct {
	UseCase *uc.UseCaseService
}

func NewKafkaHandler(uc *uc.UseCaseService) *KafkaHandler {
	return &KafkaHandler{UseCase: uc}
}