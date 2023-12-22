package nats

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Jhnvlglmlbrt/wb-order/config"
	"github.com/Jhnvlglmlbrt/wb-order/internal/models"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

type Nats struct {
	config *config.Nats
	sc     stan.Conn
	nc     *nats.Conn
}

func NewNats(cfg *config.Nats) (*Nats, error) {
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

	return &Nats{cfg, sc, nc}, nil
}

func (ns *Nats) Publish(message models.Order) error {
	order, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error at marshaling new order: %v", err)
	}

	return ns.sc.Publish(ns.config.Topic, order)
}

func (ns *Nats) Subscribe() (*models.Order, error) {

	var rec models.Order

	ch := make(chan *models.Order)

	_, err := ns.sc.Subscribe(ns.config.Topic, func(mes *stan.Msg) {

		if err := json.Unmarshal(mes.Data, &rec); err != nil {
			fmt.Printf("Error at Unmarshaling: %v\n", err)
			return
		}

		ch <- &rec
	})
	if err != nil {
		return nil, fmt.Errorf("error at subscription: %v", err)
	}

	select {
	case rec := <-ch:
		return rec, nil
	case <-time.After(60 * time.Second):
		return nil, stan.ErrTimeout
	}
}

func (ns *Nats) Close() {
	if ns.nc != nil {
		ns.nc.Close()
	}
}
