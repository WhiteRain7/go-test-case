package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	rates "github.com/WhiteRain7/GoTestCase/wallets/internal/rates/handler"
	users "github.com/WhiteRain7/GoTestCase/wallets/internal/users/handler"
	wallets "github.com/WhiteRain7/GoTestCase/wallets/internal/wallets/handler"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config struct {
	Port int `mapstructure:"port"`
}

// Status
type Status struct {
	status string
}

type Logger interface {
	Error(string)
	Warn(string)
	Info(string)
}

// status
// @Summary      Status of the service
// @Description  Returns { "status": "ok" }.
// @Tags         Other
// @Produce      json
// @Success      200 {object} Status
// @Router       /status [get]
func status(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(200)
	res.Write([]byte(`{"status":"ok"}`))
}

func NewRouter(wh *wallets.Handler, rh *rates.Handler, uh *users.Handler, swagger http.HandlerFunc) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/status", status)

	mux.HandleFunc("GET /api/v1/exchange/rates", rh.All)

	mux.HandleFunc("GET /api/v1/balance", wh.Get)
	mux.HandleFunc("POST /api/v1/wallet/deposit", wh.Deposit)
	mux.HandleFunc("POST /api/v1/wallet/withdraw", wh.Withdraw)
	mux.HandleFunc("POST /api/v1/exchange", wh.Exchange)

	mux.HandleFunc("POST /api/v1/users/register", uh.SignUp)
	mux.HandleFunc("POST /api/v1/users/login", uh.SignIn)

	mux.HandleFunc("/docs/", swagger)
	mux.Handle("/metrics", promhttp.Handler())

	return mux
}

func Run(cfg Config, router http.Handler, logger Logger) {
	srv := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port), Handler: router}

	go func() {
		srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Warn("Syscall received to shutdown server, waiting for the current tasks to be completed")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("Shutdown error, force stop: %v", err))
	}
	logger.Info("Shutdown complete")
}
