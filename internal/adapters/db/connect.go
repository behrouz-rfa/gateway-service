package db

import (
	"fmt"
	"github.com/behrouz-rfa/gateway-service/internal/core/domain"
	"github.com/behrouz-rfa/gateway-service/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	logTag = "postgresdb"
)

type DatabaseConfig struct {
	Driver string
	DSN    string
}

type DbRepository struct {
	db *gorm.DB
	lg *logger.Entry
}

func (r *DbRepository) GetDb() *gorm.DB {
	return r.db
}

func NewDB(cfg DatabaseConfig) (*DbRepository, error) {
	lg := logger.General.Component(logTag)
	var db *gorm.DB
	var err error

	switch cfg.Driver {
	case "postgres":
		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  cfg.DSN,
			PreferSimpleProtocol: true,
		}), &gorm.Config{})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(cfg.DSN), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	if err != nil {
		lg.WithError(err).Info("Error on open connection on database")
		return nil, err
	}

	err = db.AutoMigrate(&domain.User{}, &domain.Plan{})
	if err != nil {
		lg.WithError(err).Info("error on auto migration")
		return nil, err
	}

	return &DbRepository{db: db, lg: lg}, nil
}
