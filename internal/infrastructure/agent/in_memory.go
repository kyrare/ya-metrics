package agent

import "sync"

type Storage interface {
	Set(metric string, value float64)
	All() map[string]float64
	IncrementCounter()
	ResetCounter()
	GetCounter() int
}

type MemStorage struct {
	values  map[string]float64
	counter int
	mu      sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		values: make(map[string]float64),
	}
}

func (m *MemStorage) Set(metric string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.values[metric] = value
}

func (m *MemStorage) IncrementCounter() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counter++
}

func (m *MemStorage) GetCounter() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.counter
}

func (m *MemStorage) ResetCounter() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counter = 0
}

func (m *MemStorage) All() map[string]float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// копируем мапу, чтобы возвращать новую мапу, а не ссылку на нее
	// чтобы не возникала DATA RACE
	result := make(map[string]float64)
	for k, v := range m.values {
		result[k] = v
	}

	return result
}
