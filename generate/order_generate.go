package main

import (
	"fmt"
	"log"

	"github.com/Jhnvlglmlbrt/wb-order/config"
	"github.com/Jhnvlglmlbrt/wb-order/internal/cache"
	"github.com/Jhnvlglmlbrt/wb-order/internal/nats"
	"github.com/Jhnvlglmlbrt/wb-order/internal/orders/generator"
	"github.com/Jhnvlglmlbrt/wb-order/internal/repository"
	"github.com/Jhnvlglmlbrt/wb-order/package/postgres"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error at starting: %v", err)
	}

	postgresConnection, err := postgres.Connect(&cfg.PG)
	if err != nil {
		log.Fatalf("Error at Postgres connection: %v", err)
	}
	defer postgresConnection.Close()
	fmt.Println("Postgres connection successfully established.")

	repo := repository.NewRepository(postgresConnection)
	orderCache := cache.NewCache(repo)
	orderCache.Preload()

	ns, err := nats.NewNats(&cfg.Nats, orderCache)
	if err != nil {
		log.Fatalf("Error creating nats client: %v", err)
	}
	defer ns.Close()
	fmt.Println("Nats server started, connection registered")

	order := generator.GenerateOrder()
	log.Println("Order generated")

	err = order.Validate()
	if err != nil {
		fmt.Printf("Error at validating data : %v\n", err)
		return
	}

	if err := ns.Publish(*order); err != nil {
		log.Printf("Error at publishing: %v\n", err)
	} else {
		log.Println("Order sent to nats")
	}
}
