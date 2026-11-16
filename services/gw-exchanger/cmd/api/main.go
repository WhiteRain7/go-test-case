package main

import (
	"log"

	"github.com/WhiteRain7/GoTestCase/exchanger/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
