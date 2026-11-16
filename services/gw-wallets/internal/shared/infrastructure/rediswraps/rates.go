package rediswraps

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	ratesRepository "github.com/WhiteRain7/GoTestCase/wallets/internal/rates/repository"
	"github.com/redis/go-redis/v9"
)

type CachedRatesRepository struct {
	rdb      *redis.Client
	fallback *ratesRepository.Repository
	logger   Logger
}

type Logger interface {
	Info(string)
	Error(string)
}

func NewCachedRatesRepository(rdb *redis.Client, repo *ratesRepository.Repository, logger Logger) *CachedRatesRepository {
	logger.Info("New cached rates repository - Wallets")
	return &CachedRatesRepository{rdb, repo, logger}
}

func (cr *CachedRatesRepository) GetRates(ctx context.Context) (map[string]float32, error) {
	key := fmt.Sprintf("rates-%d", ctx.Value("user_id"))
	cached := cr.rdb.Get(ctx, key).Val()

	if len(cached) != 0 {
		cr.logger.Info("Cached Rates - GetRates - cache hit")
		var rates map[string]float32
		err := json.Unmarshal([]byte(cached), &rates)
		if err != nil {
			return nil, err
		}
		return rates, nil
	}

	rates, err := cr.fallback.GetRates(ctx)
	if err != nil {
		return nil, err
	}

	go func() {
		enc, err := json.Marshal(rates)
		if err == nil {
			err = cr.rdb.Set(context.Background(), key, enc, 10*time.Minute).Err()
			if err != nil {
				cr.logger.Error(fmt.Sprintf("Cached Rates - GetRates - error: %s", err.Error()))
			} else {
				cr.logger.Info("Cached Rates - GetRates - made cache")
			}
		} else {
			cr.logger.Error(fmt.Sprintf("Cached Rates - GetRates - error: %s", err.Error()))
		}
	}()

	return rates, nil
}

func (cr *CachedRatesRepository) Exchange(ctx context.Context, fromCurr string, toCurr string) (float32, error) {
	key := fmt.Sprintf("rates-%d", ctx.Value("user_id"))
	cached := cr.rdb.Get(ctx, key).Val()

	if len(cached) != 0 {
		cr.logger.Info("Cached Rates - Exchange - cache hit")
		var rates map[string]float32
		err := json.Unmarshal([]byte(cached), &rates)
		if err != nil {
			return 0, err
		}
		return rates[fromCurr] / rates[toCurr], nil
	}

	rate, err := cr.fallback.Exchange(ctx, fromCurr, toCurr)
	if err != nil {
		return 0, err
	}

	return rate, nil
}
