package broker

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

type ProdConfig struct {
	Host  string `mapstructure:"host"`
	Port  int    `mapstructure:"port"`
	Topic string `mapstructure:"topic"`
}

type BrokerProducer struct {
	kafka  *kgo.Client
	logger Logger
}

func NewProducer(cfg ProdConfig, logger Logger) (*BrokerProducer, error) {
	logger.Info(fmt.Sprintf("New producer to topic %s connected", cfg.Topic))
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)),
		kgo.DefaultProduceTopic(cfg.Topic),
	)
	if err != nil {
		return nil, err
	}
	return &BrokerProducer{cl, logger}, nil
}

func (p *BrokerProducer) SendTo(ctx context.Context, topic string, value []byte) error {
	record := &kgo.Record{Value: value, Topic: topic}
	result := p.kafka.ProduceSync(ctx, record)
	return result.FirstErr()
}

func (p *BrokerProducer) Send(ctx context.Context, value []byte) error {
	record := &kgo.Record{Value: value}
	result := p.kafka.ProduceSync(ctx, record)
	return result.FirstErr()
}

func (p *BrokerProducer) SendString(ctx context.Context, value string) error {
	return p.Send(ctx, []byte(value))
}

func (p *BrokerProducer) Close() {
	p.kafka.Close()
}
