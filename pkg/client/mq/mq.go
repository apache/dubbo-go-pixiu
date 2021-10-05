package mq

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/client/mq/impl"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/apache/dubbo-go-pixiu/pkg/filter/event"

	perrors "github.com/pkg/errors"
)

var (
	mqClient *Client
	once     sync.Once
)

func NewSingletonMQClient(config event.Config) *Client {
	if mqClient == nil {
		once.Do(func() {
			mqClient, err := NewMQClient(config)
			if err != nil {
				logger.Errorf()
			}
		})
	}
	return mqClient
}

func NewMQClient(config event.Config) (*Client, error) {
	var c *Client
	ctx, cancel := context.WithCancel(context.Background())

	switch mqType {
	case constant.MQTypeKafka:
		c = &Client{
			ctx:            ctx,
			consumerFacade: impl.NewKafkaConsumerFacade(),
			producerFacade: nil,
			stop:           cancel,
		}
	case constant.MQTypeRocketMQ:
	}
	return c, nil
}

func EventConfigToSaramaConfig(config event.Config) (sarama.Config, error) {
	c := sarama.NewConfig()
	var err error
	c.Version, err = sarama.ParseKafkaVersion(config.KafkaVersion)
	if err != nil {
		return c, err
	}
	return
}

type Client struct {
	ctx            context.Context
	consumerFacade ConsumerFacade
	producerFacade ProducerFacade
	stop           func()
}

func (c Client) Apply() error {
	panic("implement me")
}

func (c Client) Close() error {
	c.stop()
	c.consumerFacade.Stop()
	return nil
}

func (c Client) Call(req *client.Request) (res interface{}, err error) {
	body, err := ioutil.ReadAll(req.IngressRequest.Body)
	if err != nil {
		return nil, err
	}

	var mqReq event.MQRequest
	err = json.Unmarshal(body, &mqReq)
	if err != nil {
		return nil, err
	}

	paths := strings.Split(req.API.Path, "/")
	if len(paths) < 3 {
		return nil, perrors.New("failed to send message, broker or topic not found")
	}

	switch event.MQActionStrToInt[paths[0]] {
	case event.MQActionPublish:
	case event.MQActionSubscribe:
	case event.MQActionUnSubscribe:
	case event.MQActionConsumeAck:
	default:
		return nil, perrors.New("failed to get mq action")
	}

	return nil, nil
}

func (c Client) MapParams(req *client.Request) (reqData interface{}, err error) {
	return nil, perrors.New("map params does not support on mq mqClient")
}
