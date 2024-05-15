package kafka

import (
	"context"
	"encoding/json"
	"multi-site-dashboard-go/internal/config"
	"multi-site-dashboard-go/internal/domain"
	"strings"

	"github.com/segmentio/kafka-go"
)

type KafkaWriter struct {
	Writer *kafka.Writer
}

// Methods of Writer are safe to use concurrently from multiple goroutines,
// however the writer configuration should not be modified after first use.
// One writer is maintained throughout the application, and you have the 
// the ability to define the topic on a per-message basis by setting Message.Topic.
func New(cfg *config.Config) KafkaWriter {
	w := &kafka.Writer{
		Addr: kafka.TCP(strings.Split(cfg.Kafka.BrokerAddresses, ",")...),
		RequiredAcks: kafka.RequireOne,
		MaxAttempts: 5,
	}
	return KafkaWriter{Writer: w}
}

func (w KafkaWriter) PublishMachineResourceUsageEvent(ctx context.Context, arg domain.CreateMachineResourceUsageParams) error {
	v, _ := json.Marshal(arg)
	msg := kafka.Message{
		Key: []byte("testKey"),
		Value: v,
		Topic: MachineResourceUsage.String(),
	}
	if err := w.Writer.WriteMessages(ctx, msg); err != nil {
		return err
	}
	return nil
}