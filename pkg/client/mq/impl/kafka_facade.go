package impl

import (
	"context"

	"github.com/apache/dubbo-go-pixiu/pkg/client/mq"

	"github.com/Shopify/sarama"
)

func NewKafkaConsumerFacade(addrs []string, config *sarama.Config) (*KafkaConsumerFacade, error) {
	consumer, err := sarama.NewConsumer(addrs, config)
	if err != nil {
		return nil, err
	}
	return &KafkaConsumerFacade{consumer: consumer}, nil
}

type KafkaConsumerFacade struct {
	consumer sarama.Consumer
}

func (k *KafkaConsumerFacade) Subscribe(ctx context.Context, broker string, topic string, handler mq.EventHandler) error {
	panic("implement me")
}

func NewKafkaProviderFacade(addrs []string, config *sarama.Config) (*KafkaProducerFacade, error) {
	producer, err := sarama.NewSyncProducer(addrs, config)
	if err != nil {
		return nil, err
	}
	return &KafkaProducerFacade{producer: producer}, nil
}

type KafkaProducerFacade struct {
	producer sarama.SyncProducer
}

func (k *KafkaProducerFacade) Send(ctx context.Context, broker string, topic string) error {
	panic("implement me")
}
