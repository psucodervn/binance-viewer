package binance

import (
	"sync"
)

type PriceMap struct {
	m  map[string]float64
	mu sync.RWMutex
}

func NewPriceMap() *PriceMap {
	return &PriceMap{
		m: make(map[string]float64),
	}
}

func (m *PriceMap) Set(symbol string, price float64) {
	m.mu.Lock()
	m.m[symbol] = price
	m.mu.Unlock()
}

func (m *PriceMap) Get(symbol string) float64 {
	m.mu.RLock()
	p := m.m[symbol]
	m.mu.RUnlock()
	return p
}

func (m *PriceMap) GetOrStore(symbol string, price float64) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if p, ok := m.m[symbol]; ok {
		return p
	}
	m.m[symbol] = price
	return price
}
