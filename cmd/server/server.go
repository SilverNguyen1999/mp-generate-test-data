package main

import (
	"fmt"
	"log"
	"mp-generate-test-data/config"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type errorBody struct {
	Err  string `json:"error,omitempty"`
	Msg  string `json:"message,omitempty"`
	Code uint32 `json:"code,omitempty"`
	Flag uint32 `json:"flag,omitempty"`
}

func mustConnectPostgres(cfg *config.Config) *gorm.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable connect_timeout=10",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.Database)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			Colorful:      true,        // Disable color
			LogLevel:      logger.Silent,
		},
	)

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}

	err = db.Raw("SELECT 1").Error
	if err != nil {
		log.Fatal("failed to query: ", err)
	}

	pgDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	pgDB.SetMaxOpenConns(int(cfg.Postgres.MaxConnections))
	pgDB.SetMaxIdleConns(int(cfg.Postgres.IdleConnections))
	pgDB.SetConnMaxLifetime(time.Second * time.Duration(cfg.Postgres.MaxTimeLive))

	return db
}
