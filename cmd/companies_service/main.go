package main

import (
	"context"
	"flag"
	"github.com/IakimenkoD/xm-companies-service/internal/api"
	"github.com/IakimenkoD/xm-companies-service/internal/config"
	"github.com/IakimenkoD/xm-companies-service/internal/controller"
	"github.com/IakimenkoD/xm-companies-service/internal/repository/database"
	"github.com/IakimenkoD/xm-companies-service/internal/repository/dataprovider/pg"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Fatal(err.Error())
		}
	}(logger)

	var (
		configFilePath = flag.String("config", "./config.example.json", "path to configuration file")
	)

	flag.Parse()
	cfg, err := config.New(*configFilePath, logger)
	if err != nil {
		logger.Fatal("can't init config")
	}

	dbClient, err := database.NewClient(cfg)
	if err != nil {
		logger.Fatal("can't establish database connection", zap.Error(err))
	}

	if err = dbClient.Migrate(); err != nil {
		logger.Fatal("while applying database migration", zap.Error(err))

	}
	logger.Info("db migration successful")

	storage := pg.NewCompanyStorage(dbClient)
	service := controller.NewCompaniesService(cfg, storage)

	apiServer, err := api.NewServer(cfg, service)
	if err != nil {
		logger.Fatal("server init failed", zap.Error(err))
	}

	shutdown := make(chan os.Signal, 1)
	serverErrors := make(chan error, 1)

	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		serverErrors <- apiServer.ListenAndServe()
	}()

	logger.Info("service started")

	defer logger.Info("service stopped")

	select {
	case err = <-serverErrors:
		logger.Error("api server stopped", zap.Error(err))
	case sig := <-shutdown:
		logger.Info("gracefully shutdown application", zap.String("signal", sig.String()))

		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err = apiServer.Shutdown(ctx); err != nil {
			logger.Error("api server shutdown error")
			err = apiServer.Close()
		}

		if err != nil {
			logger.Error("could not stopped api server gracefully")
		}
	}
}
