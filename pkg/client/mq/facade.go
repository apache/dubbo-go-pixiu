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

// ConsumeOptions Consumer options
// TODO: Add rocketmq params
type ConsumeOptions struct {
	Topic      string
	Partition  int
	ConsumeUrl string
	CheckUrl   string
	Offset     int64
}

func (o *ConsumeOptions) ApplyOpts(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func DefaultOptions() *ConsumeOptions {
	return &ConsumeOptions{
		Topic:     "dubbo-go-pixiu-test-Topic",
		Partition: 1,
		Offset:    -2,
	}
}

type Option func(o *ConsumeOptions)

func WithTopic(t string) Option {
	return func(o *ConsumeOptions) {
		o.Topic = t
	}
}

func WithPartition(p int) Option {
	return func(o *ConsumeOptions) {
		o.Partition = p
	}
}

func WithConsumeHook(ch string) Option {
	return func(o *ConsumeOptions) {
		o.ConsumeUrl = ch
	}
}

func WithOffset(offset int64) Option {
	return func(o *ConsumeOptions) {
		o.Offset = offset
	}
}

type ProducerFacade interface {
	// Send msg to specified broker and Topic
	Send(ctx context.Context, broker string, topic string) error
}
