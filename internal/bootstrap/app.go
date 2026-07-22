package bootstrap

import (
	"errors"

	"github.com/go-api/internal/infrastructure/logger"
	"github.com/go-api/internal/infrastructure/postgres"
	"github.com/go-api/internal/shared/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type App struct {
	Config    *config.Config
	Logger    *zap.Logger
	DB        *pgxpool.Pool
	Container *Container
}

func New() (*App, error) {

	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	log, err := logger.New()
	if err != nil {
		return nil, err
	}

	db, err := postgres.NewPool(cfg.DB)
	if err != nil {
		return nil, err
	}

	container, err := NewContainer(db, cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		Config:    cfg,
		Logger:    log,
		DB:        db,
		Container: container,
	}, nil
}

func (c *Container) Close() error {
	var errs []error

	if c.Redis != nil {
		if err := c.Redis.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if c.DB != nil {
		c.DB.Close()
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
