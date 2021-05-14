package internal

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	internaldomain "github.com/MarioCarrion/todo-api/internal"
	"github.com/MarioCarrion/todo-api/internal/envvar"
)

type KafkaProducer struct {
	Producer *kafka.Producer
	Topic    string
}

// NewKafkaProducer instantiates the Kafka producer using configuration defined in environment variables.
func NewKafkaProducer(conf *envvar.Configuration) (*KafkaProducer, error) {
	host, topic, err := newKafkaConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("newKafkaConfig %w", err)
	}

	config := kafka.ConfigMap{
		"bootstrap.servers": host,
	}

	client, err := kafka.NewProducer(&config)
	if err != nil {
		return nil, fmt.Errorf("kafka.NewProducer %w", err)
	}

	return &KafkaProducer{
		Producer: client,
		Topic:    topic,
	}, nil
}

type KafkaConsumer struct {
	Consumer *kafka.Consumer
}

// NewKafkaConsumer instantiates the Kafka consumer using configuration defined in environment variables.
func NewKafkaConsumer(conf *envvar.Configuration, groupID string) (*KafkaConsumer, error) {
	host, topic, err := newKafkaConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("newKafkaConfig %w", err)
	}

	config := kafka.ConfigMap{
		"bootstrap.servers":  host,
		"group.id":           groupID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	}

	client, err := kafka.NewConsumer(&config)
	if err != nil {
		return nil, fmt.Errorf("kafka.NewConsumer %w", err)
	}

	if err := client.Subscribe(topic, nil); err != nil {
		return nil, fmt.Errorf("client.Subscribe %w", err)
	}

	return &KafkaConsumer{
		Consumer: client,
	}, nil
}

func newKafkaConfig(conf *envvar.Configuration) (host, topic string, err error) {
	host, err = conf.Get("KAFKA_HOST")
	if err != nil {
		return "", "", fmt.Errorf("conf.Get KAFKA_HOST %w", err)
	}

	topic, err = conf.Get("KAFKA_TOPIC")
	if err != nil {
		return "", "", fmt.Errorf("conf.Get KAFKA_TOPIC %w", err)
	}

	if topic == "" {
		return "", "", internaldomain.NewErrorf(internaldomain.ErrorCodeInvalidArgument, "KAFKA_TOPIC is required")
	}

	return host, topic, nil
}
