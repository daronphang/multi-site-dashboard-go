package handler

import (
	"context"
	"multi-site-dashboard-go/internal/domain"

	cv "multi-site-dashboard-go/internal/validator"

	"github.com/segmentio/kafka-go"
)

func (h *KafkaHandler) CreateMachineResourceUsageAndBroadcast(ctx context.Context, m kafka.Message) error {
	p := new(domain.CreateMachineResourceUsageParams)
	if err := cv.UnmarshalJSONAndValidate(m.Value, p); err != nil {
		return err
	}
	_, err := h.UseCase.CreateMachineResourceUsageAndBroadcast(ctx, p)
	if err != nil {
		return err
	}
	return nil
}