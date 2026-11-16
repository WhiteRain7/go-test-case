package app

import (
	"context"
	"fmt"

	"github.com/WhiteRain7/GoTestCase/analytics/internal/shared/infrastructure/clickhouseserver"
	"github.com/WhiteRain7/GoTestCase/analytics/internal/shared/infrastructure/config"
	"github.com/WhiteRain7/GoTestCase/analytics/internal/shared/infrastructure/promanalytics"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/broker"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/logger"

	transfersHandler "github.com/WhiteRain7/GoTestCase/analytics/internal/transfers/handler"
	transfersRepository "github.com/WhiteRain7/GoTestCase/analytics/internal/transfers/repository"
	transfersUsecase "github.com/WhiteRain7/GoTestCase/analytics/internal/transfers/usecase"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		return err
	}
	defer log.Close()

	consumer, err := broker.NewConsumer(cfg.Broker, log)
	if err != nil {
		return err
	}
	defer consumer.Close()
	log.Info("Kafka consumer connection established")

	ch, err := clickhouseserver.Connect(cfg.ClickHouse)
	if err != nil {
		return err
	}
	defer ch.Close()
	log.Info("Clickhouse connection established")

	transfersR := transfersRepository.NewRepository(ch, log)
	transfersUC := transfersUsecase.NewUseCase(transfersR)
	transfersH := transfersHandler.NewHandler(transfersUC)

	promanalytics.ProvideMetrics(cfg.Analytics)

	consumer.Listen(context.TODO(), func(record []byte) {
		log.Info(fmt.Sprintf("kafka record received: %s", record))

		err := transfersH.Handle(context.TODO(), record)

		if err != nil {
			log.Error(fmt.Sprintf("kafka transfer handle error: %s", err))
		}
	})

	return nil
}
