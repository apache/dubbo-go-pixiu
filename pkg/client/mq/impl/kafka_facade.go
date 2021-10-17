package impl

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/apache/dubbo-go-pixiu/pkg/client/mq"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/event"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"

	"github.com/Shopify/sarama"
	perrors "github.com/pkg/errors"
)

func NewKafkaConsumerFacade(addrs []string, config *sarama.Config) (*KafkaConsumerFacade, error) {
	consumer, err := sarama.NewConsumer(addrs, config)
	if err != nil {
		return nil, err
	}

	// does not set up cookiejar may cause some problem
	return &KafkaConsumerFacade{consumer: consumer, httpClient: &http.Client{Timeout: 5 * time.Second}}, nil
}

type KafkaConsumerFacade struct {
	consumer        sarama.Consumer
	consumerManager map[string]func()
	rwLock          sync.RWMutex
	httpClient      *http.Client
	wg              sync.WaitGroup
}

func (f *KafkaConsumerFacade) Subscribe(ctx context.Context, opts ...mq.Option) error {
	partConsumer, err := f.consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
	if err != nil {
		return err
	}
	c, cancel := context.WithCancel(ctx)
	key := mq.GetConsumerManagerKey(topic, partition)
	f.rwLock.Lock()
	defer f.rwLock.Unlock()
	f.consumerManager[key] = cancel
	f.wg.Add(2)
	go f.ConsumePartitions(c, partConsumer, consumeHook)
	go f.checkConsumeHookAlive(c, key, checkUrl)
	return nil
}

func (f *KafkaConsumerFacade) ConsumePartitions(ctx context.Context, partConsumer sarama.PartitionConsumer, consumeHook string) {
	defer f.wg.Done()

	for {
		select {
		case <-ctx.Done():
			logger.Info()
		case msg := <-partConsumer.Messages():
			data, err := json.Marshal(event.MQMsgPush{Msg: []string{string(msg.Value)}})
			if err != nil {
				logger.Warn()
				continue
			}

			req, err := http.NewRequest(http.MethodPost, consumeHook, bytes.NewReader(data))
			if err != nil {
				logger.Warn()
				continue
			}

			for i := 0; i < 5; i++ {
				err := func() error {
					resp, err := f.httpClient.Do(req)
					if err != nil {
						return err
					}
					defer resp.Body.Close()

					if resp.StatusCode == http.StatusOK {
						return nil
					}
					return perrors.New("failed send msg to consumer with status code " + strconv.Itoa(resp.StatusCode))
				}()
				if err != nil {
					logger.Warn(err.Error())
					time.Sleep(10 * time.Millisecond)
				} else {
					break
				}
			}
		}
	}
}

func (f *KafkaConsumerFacade) checkConsumeHookAlive(ctx context.Context, key string, checkUrl string) {
	defer f.wg.Done()

	ticker := time.NewTicker(15 * time.Second)
	for {
		select {
		case <-ctx.Done():
			logger.Info()
		case <-ticker.C:
			lastCheck := 0
			for i := 0; i < 5; i++ {
				err := func() error {
					req, err := http.NewRequest(http.MethodGet, checkUrl, bytes.NewReader([]byte{}))
					if err != nil {
						logger.Warn()
					}

					resp, err := f.httpClient.Do(req)
					defer resp.Body.Close()
					if err != nil {
						logger.Warn()
					}

					lastCheck = resp.StatusCode
					if resp.StatusCode != http.StatusOK {
						return perrors.New("failed check consumer alive or not with status code " + strconv.Itoa(resp.StatusCode))
					}

					logger.Warn()
					return nil
				}()
				if err != nil {
					logger.Warn()
					time.Sleep(10 * time.Millisecond)
				} else {
					break
				}
			}

			if lastCheck != http.StatusOK {
				f.consumerManager[key]()
				delete(f.consumerManager, key)
			}

		}
	}
}

func (f *KafkaConsumerFacade) UnSubscribe(opts ...mq.Option) error {
	key := mq.GetConsumerManagerKey(topic, partition)
	if cancel, ok := f.consumerManager[key]; !ok {
		return perrors.New("consumer goroutine not found")
	} else {
		cancel()
		delete(f.consumerManager, key)
	}
	return nil
}

func (f *KafkaConsumerFacade) Stop() {
	f.wg.Wait()
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

func (k *KafkaProducerFacade) Send(ctx context.Context, topic string, partition int32) error {
	panic("implement me")
}
