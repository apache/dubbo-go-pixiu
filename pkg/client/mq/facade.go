package mq

import (
	"context"
	"strconv"
	"strings"
)

type ConsumerFacade interface {
	// Subscribe message with specified broker and topic, then handle msg with handler which send msg to real consumers
	Subscribe(ctx context.Context, topic string, partition int32, consumeHook string, checkUrl string) error
	UnSubscribe(topic string, partition int32) error
	Stop()
}

func GetConsumerManagerKey(topic string, partition int32) string {
	return strings.Join([]string{topic, strconv.Itoa(int(partition))}, ".")
}

type ProducerFacade interface {
	// Send msg to specified broker and topic
	Send(ctx context.Context, broker string, topic string) error
}
