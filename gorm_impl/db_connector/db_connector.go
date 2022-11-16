package db_connector

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbClient int64

const (
	Postgres DbClient = 0
	Mysql    DbClient = 1
)

type ConnectorConfig struct {
	Mode     DbClient
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

type DatabaseConnector interface {
	GetDb() *gorm.DB
}

type databaseConnector struct {
	DatabaseConnector
	Config *ConnectorConfig
	db     *gorm.DB
}

func NewDatabaseConnector(config *ConnectorConfig) (*databaseConnector, error) {
	db, err := newClient(config)
	if err != nil || db == nil {
		return nil, err
	}
	return &databaseConnector{
		Config: config,
		db:     db,
	}, nil
}

func newClient(config *ConnectorConfig) (db *gorm.DB, err error) {
	switch config.Mode {
	case Mysql:
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Username, config.Password, config.Host, config.Port, config.Database)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Println(dsn, " - ", err)
		}
	case Postgres:
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", config.Host, config.Username, config.Password, config.Database, config.Port)
		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN: dsn,
		}), &gorm.Config{})
		if err != nil {
			fmt.Println(err)
		}
	}
	return
}

func (c *databaseConnector) GetDb() *gorm.DB {
	return c.db
}
