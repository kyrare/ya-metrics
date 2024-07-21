package connection

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

// New консруктор для создания нового соединения с БД
func New(dsn string, logger zap.SugaredLogger) (*sql.DB, error) {
	logger.Infoln("Create DB connection, dsn - ", dsn)
	return sql.Open("pgx", dsn)
}
