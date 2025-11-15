package storage

import (
	"context"

	"github.com/Roflan4eg/quiz-api/config"
	"github.com/Roflan4eg/quiz-api/pkg/logger"
	"gorm.io/gorm"
)

type SQLStorage interface {
	DB() *gorm.DB
	Close() error
}

type Container struct {
	sqlStorage SQLStorage
}

func NewContainer(cfg *config.Config) (*Container, error) {
	var (
		err error
		db  SQLStorage
	)

	logger.Debug("Initializing postgres storage")
	db, err = NewGormStorage(cfg.App.DBurl)

	if err != nil {
		return nil, err
	}

	storage := &Container{
		sqlStorage: db,
	}
	return storage, nil
}

func (c *Container) SQL() *gorm.DB {
	return c.sqlStorage.DB()
}

func (c *Container) Close(ctx context.Context) error {
	logger.Debug("Shutting down storage")
	go func() {
		err := c.sqlStorage.Close()
		if err != nil {
			logger.Warn("Error shutting down storage: ", err)
		}
	}()
	return nil
}
