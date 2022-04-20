package main

import (
	"context"
	"flag"
	"github.com/IakimenkoD/xm-companies-service/internal/api"
	"github.com/IakimenkoD/xm-companies-service/internal/config"
	"github.com/IakimenkoD/xm-companies-service/internal/controller"
	"github.com/IakimenkoD/xm-companies-service/internal/repository/database"
	"github.com/IakimenkoD/xm-companies-service/internal/repository/dataprovider/pg"
	"github.com/IakimenkoD/xm-companies-service/internal/service"
	"github.com/IakimenkoD/xm-companies-service/internal/service/http"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger, _ := zap.NewDevelopment()
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

	mq, err := service.NewMessageQueue(cfg, logger)
	if err != nil {
		logger.Fatal("while message queue init", zap.Error(err))
	}

	storage := pg.NewCompanyStorage(dbClient, logger)
	companiesService := controller.NewCompaniesService(cfg, storage, mq)
	ipChecker := http.NewIpChecker(cfg, logger)

	apiServer, err := api.NewServer(cfg, companiesService, ipChecker)
	if err != nil {
		logger.Fatal("server init failed", zap.Error(err))
	}

	shutdown := make(chan os.Signal, 1)
	serverErrors := make(chan error, 1)

	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		serverErrors <- apiServer.ListenAndServe()
	}()

	logger.Info("controller started")

	defer logger.Info("controller stopped")

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
