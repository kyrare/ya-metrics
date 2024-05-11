package metrics

import (
	"context"
	"database/sql"

	"github.com/kyrare/ya-metrics/internal/domain/metrics"
	"go.uber.org/zap"
)

type DatabaseStorage struct {
	ctx    context.Context
	db     *sql.DB
	logger zap.SugaredLogger
}

func (s *DatabaseStorage) init() error {
	_, err := s.db.ExecContext(s.ctx, "CREATE TABLE IF NOT EXISTS metrics (id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, type VARCHAR(255), name VARCHAR(255) not null, value REAL not null default 0)")

	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(s.ctx, "CREATE INDEX IF NOT EXISTS metrics_type_name ON metrics (type, name)")

	return err
}

func (s *DatabaseStorage) UpdateGauge(metric string, value float64) {
	s.update(metrics.TypeGauge, metric, value)
}

func (s *DatabaseStorage) UpdateCounter(metric string, value float64) {
	s.update(metrics.TypeCounter, metric, value)
}

func (s *DatabaseStorage) update(metricType metrics.MetricType, metric string, value float64) {
	var err error
	if s.metricExist(metricType, metric) {
		if metricType == metrics.TypeGauge {
			_, err = s.db.ExecContext(s.ctx, "UPDATE metrics SET value = $1 WHERE type = $2 AND name = $3", value, metricType, metric)
		} else {
			_, err = s.db.ExecContext(s.ctx, "UPDATE metrics SET value = value + $1 WHERE type = $2 AND name = $3", value, metricType, metric)
		}
	} else {
		_, err = s.db.ExecContext(s.ctx, "INSERT INTO metrics (type, name, value) VALUES ($1, $2, $3)", metricType, metric, value)
	}

	if err != nil {
		s.logger.Error(err)
	}
}

func (s *DatabaseStorage) GetGauges() map[string]float64 {
	return s.get(metrics.TypeGauge)
}

func (s *DatabaseStorage) GetCounters() map[string]float64 {
	return s.get(metrics.TypeCounter)
}

func (s *DatabaseStorage) get(metricType metrics.MetricType) map[string]float64 {
	result := make(map[string]float64)

	rows, err := s.db.QueryContext(s.ctx, "SELECT name, value FROM metrics WHERE type = $1", metricType)
	defer func() {
		err := rows.Close()
		if err != nil {
			s.logger.Error(err)
		}
	}()

	if err != nil {
		s.logger.Error(err)
		return result
	}

	if rows.Err() != nil {
		s.logger.Error(err)
		return result
	}

	for rows.Next() {
		var metric string
		var value float64

		err := rows.Scan(&metric, &value)
		if err != nil {
			s.logger.Error(err)
			return result
		}

		result[metric] = value
	}

	return result
}

func (s *DatabaseStorage) GetValue(metricType metrics.MetricType, metric string) (float64, bool) {
	var value float64
	row := s.db.QueryRowContext(s.ctx, "SELECT value FROM metrics WHERE type = $1 AND name = $2", metricType, metric)
	err := row.Scan(&value)

	if err == sql.ErrNoRows {
		return 0, false
	}

	return value, true
}

func (s *DatabaseStorage) Store() error {
	return nil
}

func (s *DatabaseStorage) Restore() error {
	return nil
}

func (s *DatabaseStorage) Close() error {
	return nil
}

func (s *DatabaseStorage) StoreAndClose() error {
	return nil
}

func (s *DatabaseStorage) metricExist(metricType metrics.MetricType, metric string) bool {
	row := s.db.QueryRowContext(s.ctx, "SELECT 1 FROM metrics WHERE type = $1 AND name = $2", metricType, metric)

	var exists int
	err := row.Scan(&exists)

	s.logger.Infof("Check exist %v %v, error - %v, exists - %v, result - %v", metricType, metric, err, exists, err != sql.ErrNoRows && exists == 1)

	return err != sql.ErrNoRows && exists == 1
}

func NewDatabaseStorage(ctx context.Context, DB *sql.DB, logger zap.SugaredLogger) (*DatabaseStorage, error) {
	s := &DatabaseStorage{
		ctx:    ctx,
		db:     DB,
		logger: logger,
	}

	if err := s.init(); err != nil {
		return nil, err
	}

	return s, nil
}
