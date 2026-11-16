package broker

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

type ConConfig struct {
	Host  string `mapstructure:"host"`
	Port  int    `mapstructure:"port"`
	Group string `mapstructure:"group"`
	Topic string `mapstructure:"topic"`
}

type BrokerConsumer struct {
	kafka  *kgo.Client
	logger Logger
}

type Logger interface {
	Info(string)
	Error(string)
}

func NewConsumer(cfg ConConfig, logger Logger) (*BrokerConsumer, error) {
	logger.Info(fmt.Sprintf("New consumer on topic %s with group %s connected", cfg.Topic, cfg.Group))
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)),
		kgo.ConsumerGroup(cfg.Group),
		kgo.ConsumeTopics(cfg.Topic),
	)
	if err != nil {
		return nil, err
	}
	return &BrokerConsumer{cl, logger}, nil
}

func (c *BrokerConsumer) Listen(ctx context.Context, handler func([]byte)) {
	for {
		fetches := c.kafka.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				c.logger.Error(fmt.Sprintf("error: %s on topic %s", err.Err.Error(), err.Topic))
			}
		}
		fetches.EachRecord(func(record *kgo.Record) {
			handler(record.Value)
		})
	}
}

func (c *BrokerConsumer) Close() {
	c.kafka.Close()
}
