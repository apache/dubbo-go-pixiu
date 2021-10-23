package mq

import (
	"context"
	"strconv"
	"strings"
)

type ConsumerFacade interface {
	// Subscribe message with specified broker and Topic, then handle msg with handler which send msg to real consumers
	Subscribe(ctx context.Context, option ...Option) error
	UnSubscribe(opts ...Option) error
	Stop()
}

func GetConsumerManagerKey(topic string, partition int32) string {
	return strings.Join([]string{topic, strconv.Itoa(int(partition))}, ".")
}

// MQOptions Consumer options
// TODO: Add rocketmq params
type MQOptions struct {
	Topic      string
	Partition  int
	ConsumeUrl string
	CheckUrl   string
	Offset     int64
}

func (o *MQOptions) ApplyOpts(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func DefaultOptions() *MQOptions {
	return &MQOptions{
		Topic:     "dubbo-go-pixiu-default-topic",
		Partition: 1,
		Offset:    -2,
	}
}

type Option func(o *MQOptions)

func WithTopic(t string) Option {
	return func(o *MQOptions) {
		o.Topic = t
	}
}

func WithPartition(p int) Option {
	return func(o *MQOptions) {
		o.Partition = p
	}
}

func WithConsumeUrl(ch string) Option {
	return func(o *MQOptions) {
		o.ConsumeUrl = ch
	}
}

func WithOffset(offset int64) Option {
	return func(o *MQOptions) {
		o.Offset = offset
	}
}

type ProducerFacade interface {
	// Send msg to specified broker and Topic
	Send(msgs []string, opts ...Option) error
}
