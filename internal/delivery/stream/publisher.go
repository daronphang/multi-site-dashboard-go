package stream

import (
	"context"
	"encoding/json"
	"multi-site-dashboard-go/internal/config"
	"multi-site-dashboard-go/internal/domain"
	"strings"

	"github.com/segmentio/kafka-go"
)

type KafkaWriter struct {
	writer *kafka.Writer
}

// Methods of Writer are safe to use concurrently from multiple goroutines,
// however the writer configuration should not be modified after first use.
// One writer is maintained throughout the application, and you have the 
// the ability to define the topic on a per-message basis by setting Message.Topic.
func New(cfg *config.Config) KafkaWriter {
	w := &kafka.Writer{
		Addr: kafka.TCP(strings.Split(cfg.Kafka.BrokerAddresses, ",")...),
		RequiredAcks: kafka.RequireOne,
	}
	return KafkaWriter{writer: w}
}

func (w KafkaWriter) PublishToMachineResourceUsage(ctx context.Context, arg domain.CreateMachineResourceUsageParams) error {
	v, _ := json.Marshal(arg)
	msg := kafka.Message{
		Key: []byte("testKey"),
		Value: v,
		Topic: MachineResourceUsage.String(),
	}
	// TODO: retry message on failure with exponential backoff
	if err := w.writer.WriteMessages(ctx, msg); err != nil {
		return err
	}
	return nil
}