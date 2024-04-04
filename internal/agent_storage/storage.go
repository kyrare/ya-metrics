package agent_storage

// todo По заданию не понял, что должен делать агент с данными, пока завел отдельный Storage

type Storage interface {
	Set(metric string, value float64)
	Increment(metric string)
	All() map[string]float64
}

type MemStorage struct {
	Values map[string]float64
}

func (m *MemStorage) Set(metric string, value float64) {
	m.Values[metric] = value
}

func (m *MemStorage) Increment(metric string) {
	if v, ok := m.Values[metric]; ok {
		m.Values[metric] = v + 1
	} else {
		m.Values[metric] = 1
	}
}

func (m *MemStorage) All() map[string]float64 {
	return m.Values
}
