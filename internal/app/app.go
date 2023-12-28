package app

import (
	"fmt"
	"log"

	"github.com/Jhnvlglmlbrt/wb-order/config"
	"github.com/Jhnvlglmlbrt/wb-order/internal/cache"
	"github.com/Jhnvlglmlbrt/wb-order/internal/nats"
	"github.com/Jhnvlglmlbrt/wb-order/internal/orders/controller"
	"github.com/Jhnvlglmlbrt/wb-order/internal/orders/generator"
	"github.com/Jhnvlglmlbrt/wb-order/internal/repository"
	"github.com/Jhnvlglmlbrt/wb-order/package/http"
	"github.com/Jhnvlglmlbrt/wb-order/package/postgres"
)

func Run(cfg *config.Config) {

	postgresConnection, err := postgres.Connect(&cfg.PG)
	if err != nil {
		log.Fatalf("Error at Postgres connection: %v", err)
	}
	defer postgresConnection.Close()
	fmt.Println("Postgres connection successfully established.")

	repo := repository.NewRepository(postgresConnection)

	creationTableError := repo.CreateTable()
	if creationTableError != nil {
		log.Fatalf("Error at table creation: %v", err)
	}
	fmt.Println("Table created.")

	orderCache := cache.NewCache(repo)
	orderCache.Preload()

	ns, err := nats.NewNats(&cfg.Nats, orderCache)
	if err != nil {
		log.Fatalf("Error creating nats client: %v", err)
	}
	defer ns.Close()
	fmt.Println("Nats server started, connection registered")

	// go func() {
	// 	orderGeneratorTicker := time.NewTicker(30 * time.Second)
	// 	defer orderGeneratorTicker.Stop()

	// 	for range orderGeneratorTicker.C {
	// 		order := generator.GenerateOrder()
	// 		log.Println("Order generated")
	// 		if err := ns.Publish(*order); err != nil {
	// 			log.Printf("Error at publishing: %v\n", err)
	// 		}
	// 		log.Println("Order sent to nats")
	// 	}
	// }()

	go func() {
		order := generator.GenerateOrder()
		log.Println("Order generated")
		if err := ns.Publish(*order); err != nil {
			log.Printf("Error at publishing: %v\n", err)
		}
		log.Println("Order sent to nats")
	}()

	go func() {
		if ns != nil {
			if err := ns.InitSubscription(); err != nil {
				log.Printf("Error initializing subscription: %v\n", err)
			}
		} else {
			log.Println("Nats client is nil. Skipping subscription.")
		}
	}()

	// go func() {
	// 	orderSubscriberTicker := time.NewTicker(15 * time.Second)
	// 	defer orderSubscriberTicker.Stop()

	// 	for range orderSubscriberTicker.C {
	// 		if ns != nil {
	// 			order, err := ns.Subscribe() // TODO: сделать что то с durable подпиской - ошибка "duplicate durable registration"
	// 			if err != nil {
	// 				log.Printf("Error subscribing: %v\n", err)
	// 			} else {
	// 				log.Printf("Received order: %+v\n", order)
	// 				orderCache.Init(*order)
	// 			}
	// 		} else {
	// 			log.Println("Nats client is nil. Skipping subscription.")
	// 		}
	// 	}
	// }()

	// go func() {
	// 	for {
	// 		if ns != nil {
	// 			order, err := ns.Subscribe()
	// 			if err != nil {
	// 				log.Printf("Error subscribing: %v\n", err)
	// 			} else {
	// 				log.Printf("Received order: %+v\n", order)
	// 				orderCache.Init(*order)
	// 			}
	// 		} else {
	// 			log.Println("Nats client is nil. Skipping subscription.")
	// 		}
	// 	}
	// }()

	// subscribe and save order in db and in cache

	httpServer := http.NewServer(&cfg.HTTP)
	orderController := controller.NewController(orderCache)
	servertStart := httpServer.Start(orderController.HandlePage, orderController.HandleGetData)
	if servertStart != nil {
		log.Fatalf("Error at server starting: %v", servertStart)
	}
	// get order, all orders, create order
	// server start
}
