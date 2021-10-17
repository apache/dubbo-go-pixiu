package mq

import (
	"context"
	"strconv"
	"strings"
)

type ConsumerFacade interface {
	// Subscribe message with specified broker and topic, then handle msg with handler which send msg to real consumers
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
	topic       string
	partition   int
	consumeHook string
	checkUrl    string
	offset      int64
}

func DefaultOptions() *ConsumeOptions {
	return &ConsumeOptions{
		topic:     "dubbo-go-pixiu-test-topic",
		partition: 1,
		offset:    -2,
	}
}

type Option func(o *ConsumeOptions)

func WithTopic(t string) Option {
	return func(o *ConsumeOptions) {
		o.topic = t
	}
}

func WithPartition(p int) Option {
	return func(o *ConsumeOptions) {
		o.partition = p
	}
}

func WithConsumeHook(ch string) Option {
	return func(o *ConsumeOptions) {
		o.consumeHook = ch
	}
}

func WithOffset(offset int64) Option {
	return func(o *ConsumeOptions) {
		o.offset = offset
	}
}

type ProducerFacade interface {
	// Send msg to specified broker and topic
	Send(ctx context.Context, broker string, topic string) error
}
