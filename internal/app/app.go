package app

import (
	"fmt"
	"log"

	"github.com/Jhnvlglmlbrt/wb-order/config"
	"github.com/Jhnvlglmlbrt/wb-order/internal/nats"
	"github.com/Jhnvlglmlbrt/wb-order/internal/repository"
	"github.com/Jhnvlglmlbrt/wb-order/package/postgres"
)

func Run(cfg *config.Config) {

	ns, err := nats.NewNats(&cfg.Nats)
	if err != nil {
		log.Fatalf("Error creating nats client: %v", err)
	}
	defer ns.Close()

	fmt.Println("Nats server started, connection registered")

	postgresConnection, err := postgres.Connect(&cfg.PG)
	if err != nil {
		log.Fatalf("Error at Postgres connection: %v", err)
	}
	defer postgresConnection.Close()

	repo := repository.NewRepository(postgresConnection)

	creationTableError := repo.CreateTable()
	if creationTableError != nil {
		log.Fatalf("Error at table creation: %v", err)
	}

	// Cache initialization
	// preload cache

	// publish info

	// subscribe and save order in db and in cache

	// server init
	// controller init
	// get order, all orders, create order
	// server start
}
