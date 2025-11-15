package storage

import (
	"fmt"

	"github.com/Roflan4eg/quiz-api/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormStorage struct {
	db *gorm.DB
}

func NewGormStorage(databaseURL string) (*GormStorage, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established")
	return &GormStorage{db: db}, nil
}

func (s *GormStorage) DB() *gorm.DB {
	return s.db
}

func (s *GormStorage) Close() error {
	if s.db == nil {
		return nil
	}
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	return sqlDB.Close()
}

func (s *GormStorage) HealthCheck() error {
	if s.db == nil {
		return fmt.Errorf("database connection is nil")
	}
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	return sqlDB.Ping()
}
