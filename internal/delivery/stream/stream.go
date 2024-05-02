package stream

import (
	"context"
	"fmt"
	"multi-site-dashboard-go/internal"
	"multi-site-dashboard-go/internal/config"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

var logger, _ = internal.WireLogger()

type Topic int32

const (
	MachineResourceUsage Topic = iota
)

func (t Topic) String() string {
	switch t {
	case MachineResourceUsage:
		return "machine-resource-usage"
	}
	return "unknown"
}

func CreateKafkaTopics(cfg *config.Config) error {
	conn, err := kafka.Dial("tcp", strings.Split(cfg.Kafka.BrokerAddresses, ",")[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	/* 
	Config for each topic must be explicitly set for the following reasons:
	- You cannot decrease the number of partitions 
	- Increasing the partitions will force a rebalance
	- ReplicationFactor cannot be greater than the number of brokers available
	*/
	topicConfigs := []kafka.TopicConfig{
		{
			Topic:  MachineResourceUsage.String(),
			NumPartitions: 10,
			ReplicationFactor: 1,
		},
	}

	if err := controllerConn.CreateTopics(topicConfigs...); err != nil {
		return err
	}

	return nil
}

func GracefulShutdown(ctx context.Context, w KafkaWriter) {
	fmt.Printf("performing graceful shutdown with timeout of %v...", 10*time.Second)
	if err := w.writer.Close(); err != nil {
		logger.Error("failed to close Kafka writer", zap.String("trace", err.Error()))
	}
}