package goutils

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresStore implements HistoricalStore backed by PostgreSQL using GORM.
type PostgresStore struct {
	mu     sync.RWMutex
	db     *gorm.DB
	logger *logrus.Logger
}

func NewPostgresStore(dsn string, logger *logrus.Logger) (*PostgresStore, error) {
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}
	return &PostgresStore{db: db, logger: logger}, nil
}

// Reconnect closes the current connection and opens a new one with the given DSN.
func (s *PostgresStore) Reconnect(dsn string) error {
	newDB, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return err
	}
	s.mu.Lock()
	oldDB := s.db
	s.db = newDB
	s.mu.Unlock()
	if sqlDB, err := oldDB.DB(); err == nil {
		err = sqlDB.Close()
		if err != nil {
			s.logger.Warn("Failed to close old database connection", "error", err)
			return fmt.Errorf("failed to close old database connection: %w", err)
		}
	}
	return nil
}

func (s *PostgresStore) conn() *gorm.DB {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.db
}

func (s *PostgresStore) DB() *gorm.DB {
	return s.conn()
}

func (s *PostgresStore) Close() error {
	sqlDB, err := s.conn().DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
