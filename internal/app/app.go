package app

import (
	"github.com/WhiteRain7/GoTestCase/internal/shared/infrastructure/config"
	httpserver "github.com/WhiteRain7/GoTestCase/internal/shared/infrastructure/http"
	"github.com/WhiteRain7/GoTestCase/internal/shared/infrastructure/logger"
	"github.com/WhiteRain7/GoTestCase/internal/shared/infrastructure/postgres"
	walletsHandler "github.com/WhiteRain7/GoTestCase/internal/wallets/handler"
	walletsRepository "github.com/WhiteRain7/GoTestCase/internal/wallets/repository"
	walletsUsecase "github.com/WhiteRain7/GoTestCase/internal/wallets/usecase"
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

	db, err := postgres.Connect(cfg.DB)
	if err != nil {
		return err
	}
	defer db.Close()
	log.Info("DB connected")

	walletsRepo := walletsRepository.NewRepository(db, log)
	walletsUC := walletsUsecase.NewUseCase(walletsRepo, log)
	walletsH := walletsHandler.NewHandler(walletsUC, log)

	router := httpserver.NewRouter(walletsH)

	log.Info("Run application")
	httpserver.Run(cfg.HTTP, router)
	return nil
}
