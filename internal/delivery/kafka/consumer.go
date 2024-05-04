package kafka

import (
	"context"
	"strings"

	"multi-site-dashboard-go/internal/config"
	"multi-site-dashboard-go/internal/delivery/kafka/handler"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// To execute the consumer functions as a goroutine.
// One consumer per thread/goroutine is the rule.
// Creating more consumers than the number of partitions will result in unused consumers.

func ConsumeMsgFromMachineResourceUsageTopic(ctx context.Context, cfg *config.Config, h *handler.KafkaHandler) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: strings.Split(cfg.Kafka.BrokerAddresses, ","),
		GroupID: MachineResourceUsage.String(),
		Topic: MachineResourceUsage.String(),
	})
	for {
		select {
		case <- ctx.Done():
			// Important to close reader when process exits.
			if err := r.Close(); err != nil {
				logger.Error(
					"unable to close reader for MachineResourceUsage topic",
					zap.String("trace", err.Error()),
				)
			}
			return
		default:
			m, err := r.ReadMessage(ctx)
			if err != nil {
				logger.Error(
					"error reading message from MachineResourceUsage topic",
					zap.String("trace", err.Error()),
				)
				break
			}
			logger.Info(
				"message received",
				zap.Int64("offset", m.Offset),
				zap.String("key", string(m.Key)),
				zap.String("value", string(m.Value)),
			)
			if err := h.CreateMachineResourceUsageAndBroadcast(ctx, m); err != nil {
				logger.Error(
					"error processing message from MachineResourceUsage topic",
					zap.String("trace", err.Error()),
				)
			}
		}
	}
}