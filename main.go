package main

import (
	"log"

	"github.com/Jhnvlglmlbrt/wb-order/config"
	"github.com/Jhnvlglmlbrt/wb-order/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error at starting: %v", err)
	}

	app.Run(cfg)
}
