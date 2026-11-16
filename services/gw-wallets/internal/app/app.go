package app

import (
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/broker"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/logger"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/postgres"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/protoclient"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/redis"
	codederrors "github.com/WhiteRain7/GoTestCase/pkg/utils/coded-errors"
	ratesHandler "github.com/WhiteRain7/GoTestCase/wallets/internal/rates/handler"
	ratesRepository "github.com/WhiteRain7/GoTestCase/wallets/internal/rates/repository"
	ratesUsecase "github.com/WhiteRain7/GoTestCase/wallets/internal/rates/usecase"
	authorization "github.com/WhiteRain7/GoTestCase/wallets/internal/shared/infrastructure/auth"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/shared/infrastructure/config"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/shared/infrastructure/httpserver"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/shared/infrastructure/promanalytics"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/shared/infrastructure/rediswraps"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/shared/infrastructure/swagger"
	walletsHandler "github.com/WhiteRain7/GoTestCase/wallets/internal/wallets/handler"
	walletsRepository "github.com/WhiteRain7/GoTestCase/wallets/internal/wallets/repository"
	walletsUsecase "github.com/WhiteRain7/GoTestCase/wallets/internal/wallets/usecase"
	"github.com/WhiteRain7/GoTestCase/wallets/pkg/txmanager"

	usersHandler "github.com/WhiteRain7/GoTestCase/wallets/internal/users/handler"
	usersRepository "github.com/WhiteRain7/GoTestCase/wallets/internal/users/repository"
	usersUsecase "github.com/WhiteRain7/GoTestCase/wallets/internal/users/usecase"
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

	auth := authorization.MustLoad(cfg.Auth)
	log.Info("JWT auth module loaded")

	db, err := postgres.Connect(cfg.DB)
	if err != nil {
		return err
	}
	defer db.Close()
	log.Info("DB connected")

	dbx := txmanager.NewTXManager(db)

	exchanger, exchangerConn, err := protoclient.Connect(cfg.Proto)
	if err != nil {
		return err
	}
	defer exchangerConn.Close()
	log.Info("gRPC connection with Exchanger service established")

	rdb, err := redis.Connect(cfg.Redis)
	if err != nil {
		return err
	}
	log.Info("redis connection established")

	producer, err := broker.NewProducer(cfg.Broker, log)
	if err != nil {
		return err
	}
	defer producer.Close()
	log.Info("kafka producer connection established")

	promanalytics.Init()

	codederrors.SetDebug(cfg.Debug)

	ratesRepo := ratesRepository.NewRepository(exchanger, log)
	cachedRatesRepo := rediswraps.NewCachedRatesRepository(rdb, ratesRepo, log)
	ratesUC := ratesUsecase.NewUseCase(cachedRatesRepo, log)
	ratesH := ratesHandler.NewHandler(ratesUC, log)

	walletsRepo := walletsRepository.NewRepository(dbx, log)
	walletsUC := walletsUsecase.NewUseCase(dbx, walletsRepo, cachedRatesRepo, producer, log)
	walletsH := walletsHandler.NewHandler(walletsUC, log)

	usersRepo := usersRepository.NewRepository(dbx, log)
	usersUC := usersUsecase.NewUseCase(dbx, usersRepo, walletsRepo, log)
	usersH := usersHandler.NewHandler(usersUC, auth, log)

	swaggerui, err := swagger.Load(cfg.Swagger)
	if err != nil {
		return err
	}

	router := httpserver.NewRouter(walletsH, ratesH, usersH, swaggerui)

	router = promanalytics.MiddlewareMetrics(router)
	router = log.MiddlewareLogRequests(router)
	router = auth.MiddlewareHandleJWT(router)

	log.Info("Run application")
	httpserver.Run(cfg.HTTP, router, log)
	return nil
}
