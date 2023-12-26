package cache

import (
	"fmt"
	"log"
	"sync"

	"github.com/Jhnvlglmlbrt/wb-order/internal/models"
	"github.com/Jhnvlglmlbrt/wb-order/internal/repository"
)

type Cache interface {
	Init(or models.Order)
	Preload()
	GetOrderByUid(uid string) *models.Order
	GetOrders() map[string]*models.Order
}

type cache struct {
	cache map[string]*models.Order
	repo  repository.Repository
	mu    sync.Mutex
}

func NewCache(repo repository.Repository) Cache {
	return &cache{
		cache: map[string]*models.Order{},
		repo:  repo,
		mu:    sync.Mutex{},
	}
}

func (c *cache) Init(or models.Order) {
	if err := c.repo.SaveOrder(or); err != nil {
		log.Printf("cannot insert order: %v", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[or.OrderUid] = &or

	log.Printf("Cache written: %s\n", or.OrderUid)
}

func (c *cache) Preload() {
	ords, err := c.repo.GetAll()
	if err != nil {
		fmt.Printf("Error at getting all data: %v\n", err)
	}

	fmt.Printf("DB len: %d\n", len(ords))

	c.mu.Lock()
	defer c.mu.Unlock()
	for _, or := range ords {
		c.cache[or.OrderUid] = &or
	}
}

func (c *cache) GetOrderByUid(uid string) *models.Order {
	return c.cache[uid]
}

func (c *cache) GetOrders() map[string]*models.Order {
	return c.cache
}
