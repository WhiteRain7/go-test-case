package app

import (
	"context"
	"fmt"

	"github.com/WhiteRain7/GoTestCase/notification/internal/shared/infrastructure/config"
	"github.com/WhiteRain7/GoTestCase/notification/internal/shared/infrastructure/mongoserver"
	"github.com/WhiteRain7/GoTestCase/notification/internal/shared/infrastructure/promanalytics"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/broker"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/logger"

	transfersHandler "github.com/WhiteRain7/GoTestCase/notification/internal/transfers/handler"
	transfersRepository "github.com/WhiteRain7/GoTestCase/notification/internal/transfers/repository"
	transfersUsecase "github.com/WhiteRain7/GoTestCase/notification/internal/transfers/usecase"
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
	log.Info("kafka consumer connection established")

	mdb, err := mongoserver.Connect(cfg.Mongo)
	if err != nil {
		return err
	}
	defer mdb.Disconnect(context.Background())
	log.Info("Mongo connection established")

	transfersR := transfersRepository.NewRepository(mdb)
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
