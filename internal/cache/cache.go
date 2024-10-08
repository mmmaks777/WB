package cache

import (
	"sync"
	"wb_l0/internal/domain"

	lru "github.com/hashicorp/golang-lru"
)

type OrderCache struct {
	cache *lru.Cache
	mu    sync.RWMutex
}

func NewOrderCache(size int) (*OrderCache, error) {
	c, err := lru.New(size)
	if err != nil {
		return nil, err
	}
	return &OrderCache{cache: c}, nil
}

func (o *OrderCache) Set(order domain.Order) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.cache.Add(order.OrderUID, order)
}

func (o *OrderCache) Get(orderUID string) (domain.Order, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if value, ok := o.cache.Get(orderUID); ok {
		return value.(domain.Order), true
	}
	return domain.Order{}, false
}
