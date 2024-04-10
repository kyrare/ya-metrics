package agent

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
}

func (m *MemStorage) Set(metric string, value float64) {
	m.values[metric] = value
}

func (m *MemStorage) IncrementCounter() {
	m.counter++
}

func (m *MemStorage) GetCounter() int {
	return m.counter
}

func (m *MemStorage) ResetCounter() {
	m.counter = 0
}

func (m *MemStorage) All() map[string]float64 {
	return m.values
}

func NewMemeStorage() *MemStorage {
	return &MemStorage{
		values: make(map[string]float64),
	}
}
