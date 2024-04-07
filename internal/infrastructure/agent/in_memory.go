package agent

type Storage interface {
	Set(metric string, value float64)
	Increment(metric string)
	All() map[string]float64
}

type MemStorage struct {
	values map[string]float64
}

func (m *MemStorage) Set(metric string, value float64) {
	m.values[metric] = value
}

func (m *MemStorage) Increment(metric string) {
	if v, ok := m.values[metric]; ok {
		m.values[metric] = v + 1
	} else {
		m.values[metric] = 1
	}
}

func (m *MemStorage) All() map[string]float64 {
	return m.values
}

func NewMemeStorage() *MemStorage {
	return &MemStorage{
		values: make(map[string]float64),
	}
}
