package main

import (
	"fmt"
	"githib.com/dkischenko/company-api/configs"
	"githib.com/dkischenko/company-api/internal/app"
	"githib.com/dkischenko/company-api/internal/company"
	"githib.com/dkischenko/company-api/internal/company/database"
	"githib.com/dkischenko/company-api/models"
	"githib.com/dkischenko/company-api/pkg/logger"
	"github.com/caarlos0/env"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%s\n", err)
	}
}

func run() error {
	l, err := logger.GetLogger()
	if err != nil {
		return fmt.Errorf("cannot init logger: %w", err)
	}
	router := mux.NewRouter()
	cfg := configs.Config{}
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("cannot parse config file: %w", err)
	}

	//connect to DB
	db, err := gorm.Open(postgres.Open(cfg.DatabaseDsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("cannot connect to database: %w", err)
	}

	err = db.AutoMigrate(models.Company{}, models.User{})
	if err != nil {
		return fmt.Errorf("cannot migrate database: %w", err)
	}

	accessTokenTTL, err := time.ParseDuration(cfg.AccessTokenTTL)
	if err != nil {
		return fmt.Errorf("cannot parse token: %w", err)
	}

	storage := database.NewStorage(db, l)
	service := company.NewService(l, storage, accessTokenTTL)
	handler := company.NewHandler(l, service, &cfg)
	handler.Register(router)
	app.RunServer(router, l, &cfg)

	return nil
}
