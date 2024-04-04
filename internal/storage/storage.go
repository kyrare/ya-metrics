package storage

type Storage interface {
	UpdateGauge(metric string, value float64)
	UpdateCounter(metric string, value float64)
}

type MemStorage struct {
	Gauges   map[string]float64
	Counters map[string][]float64
}

func (storage *MemStorage) UpdateGauge(metric string, value float64) {
	if storage.Gauges == nil {
		storage.Gauges = make(map[string]float64)
	}

	storage.Gauges[metric] = value
}

func (storage *MemStorage) UpdateCounter(metric string, value float64) {
	if storage.Counters == nil {
		storage.Counters = make(map[string][]float64)
	}

	storage.Counters[metric] = append(storage.Counters[metric], value)
}
