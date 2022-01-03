package main

import (
	"github.com/Baraulia/goLab/IndTask.git/internal/config"
	"github.com/Baraulia/goLab/IndTask.git/internal/handler"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
	"github.com/Baraulia/goLab/IndTask.git/internal/service"
	"github.com/Baraulia/goLab/IndTask.git/pkg/logging"
	"github.com/Baraulia/goLab/IndTask.git/pkg/postgres"
	"github.com/Baraulia/goLab/IndTask.git/pkg/server"
	_ "github.com/lib/pq"
)

func main() {
	logger := logging.GetLogger()
	cfg := config.GetConfig()

	db, err := postgres.NewPostgresDB(postgres.PostgresDB{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBName,
		SSLMode:  cfg.DB.SSLMode,
	})
	if err != nil {
		logger.Panicf("failed to initialize db:%s", err.Error())
	}

	Rep := repository.NewRepository(db)
	services := service.NewService(Rep)
	handlers := handler.NewHandler(logger, services)
	srv := new(server.Server)
	logger.Infof("Running server on %s:%s...", cfg.Listen.BindIp, cfg.Listen.Port)
	if err := srv.Run(cfg.Listen.BindIp, cfg.Listen.Port, handlers.InitRoutes()); err != nil {
		logger.Panicf("Error occured while running http server: %s", err.Error())
	}

}
