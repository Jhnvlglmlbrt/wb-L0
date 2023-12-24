package cache

import (
	"fmt"
	"sync"

	"github.com/Jhnvlglmlbrt/wb-order/internal/models"
	"github.com/Jhnvlglmlbrt/wb-order/internal/repository"
)

type Cache struct {
	cache map[string]*models.Order
	repo  repository.Repository
	mu    sync.Mutex
}

func NewCache(repo repository.Repository) *Cache {
	return &Cache{
		cache: map[string]*models.Order{},
		repo:  repo,
		mu:    sync.Mutex{},
	}
}

func (c *Cache) CacheInit(or models.Order) {
	if err := c.repo.SaveOrder(or); err != nil {
		fmt.Printf("cannot insert order: %v", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[or.OrderUid] = &or

	fmt.Printf("Cache written: %s\n", or.OrderUid)
}

func (c *Cache) Preload() {
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

func (c *Cache) GetOrderByUid(uid string) *models.Order {
	return c.cache[uid]
}

func (c *Cache) GetOrders() map[string]*models.Order {
	return c.cache
}
