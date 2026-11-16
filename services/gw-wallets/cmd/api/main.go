package main

import (
	"log"

	"github.com/WhiteRain7/GoTestCase/wallets/internal/app"
)

// @title Wallets API
// @version 1.0.0
// @host localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Write "Bearer [token]", with token got from /login.
func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
