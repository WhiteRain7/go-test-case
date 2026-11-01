package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/WhiteRain7/GoTestCase/internal/shared/infrastructure/config"
	wallets "github.com/WhiteRain7/GoTestCase/internal/wallets/handler"
)

func NewRouter(wh *wallets.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/status", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		if v, err := json.Marshal(struct{ Ok bool }{true}); err == nil {
			res.Write(v)
		}
	})

	mux.HandleFunc("GET /api/v1/wallets/{id}", wh.Get)
	mux.HandleFunc("POST /api/v1/wallets/create", wh.Create)
	mux.HandleFunc("POST /api/v1/wallet", wh.Update)

	return mux
}

func Run(ctx config.HTTPConfig, router *http.ServeMux) {
	http.ListenAndServe(fmt.Sprintf(":%v", ctx.Port), router)
}
