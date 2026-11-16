package app

import (
	"github.com/WhiteRain7/GoTestCase/exchanger/internal/adapters/grpc"
	ratesRepository "github.com/WhiteRain7/GoTestCase/exchanger/internal/rates/repository"
	ratesUsecase "github.com/WhiteRain7/GoTestCase/exchanger/internal/rates/usecase"
	"github.com/WhiteRain7/GoTestCase/exchanger/internal/shared/infrastructure/config"
	"github.com/WhiteRain7/GoTestCase/exchanger/internal/shared/infrastructure/promanalytics"
	"github.com/WhiteRain7/GoTestCase/exchanger/internal/shared/infrastructure/protoserver"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/logger"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/postgres"
)

func Run() error {
	cfg, err := config.MustLoad()
	if err != nil {
		return err
	}

	log, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		return err
	}
	defer log.Close()

	db, err := postgres.Connect(cfg.DB)
	if err != nil {
		return err
	}
	defer db.Close()
	log.Info("DB connected")

	ratesRepo := ratesRepository.NewRepository(db, log)
	ratesUC := ratesUsecase.NewUseCase(ratesRepo, log)
	exchangeService := grpc.NewExchangeService(ratesUC)

	promanalytics.ProvideMetrics(cfg.Analytics)

	server := protoserver.NewServer(
		exchangeService,
	)

	protoserver.Run(cfg.Proto, server)

	log.Info("Run application")
	return nil
}
