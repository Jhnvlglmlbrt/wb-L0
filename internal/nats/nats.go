package nats

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Jhnvlglmlbrt/wb-order/config"
	"github.com/Jhnvlglmlbrt/wb-order/internal/cache"
	"github.com/Jhnvlglmlbrt/wb-order/internal/models"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

type Nats struct {
	config     *config.Nats
	sc         stan.Conn
	nc         *nats.Conn
	sub        stan.Subscription
	orderCache *cache.Cache
}

func NewNats(cfg *config.Nats, orderCache *cache.Cache) (*Nats, error) {
	natsUrl := fmt.Sprintf("nats://%s:%s", cfg.Host, cfg.Port)

	nc, err := nats.Connect(natsUrl)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to nats: %v", err)
	}

	sc, err := stan.Connect(cfg.Cluster, cfg.Client,
		stan.Pings(10, 3),
		stan.SetConnectionLostHandler(func(_ stan.Conn, err error) {
			fmt.Printf("Connection lost: %v", err)
		}))
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("cannot connect to stan: %v", err)
	}

	return &Nats{cfg, sc, nc, nil, orderCache}, nil
}

func (ns *Nats) Publish(message models.Order) error {
	order, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error at marshaling new order: %v", err)
	}

	return ns.sc.Publish(ns.config.Topic, order)
}

func (ns *Nats) InitSubscription() error {
	var rec models.Order

	if ns.sub == nil {
		sub, err := ns.sc.Subscribe(ns.config.Topic, func(mes *stan.Msg) {
			if err := json.Unmarshal(mes.Data, &rec); err != nil {
				fmt.Printf("Error at Unmarshaling: %v\n", err)
				return
			}

			ns.handleMessage(&rec)
		}, stan.DurableName(ns.config.Durable))
		if err != nil {
			return fmt.Errorf("error at subscription: %v", err)
		}
		ns.sub = sub
	}

	return nil
}

func (ns *Nats) handleMessage(order *models.Order) {
	log.Printf("Received order: %+v\n", order)
	ns.orderCache.Init(*order)
}

func (ns *Nats) Close() {
	if ns.nc != nil {
		ns.nc.Close()
	}
	if ns.sub != nil {
		ns.sub.Close()
	}
}
