package main

import (
	"context"
	"fmt"

	"github.com/calvinbui/teleport-sre-challenge/config"
	"github.com/calvinbui/teleport-sre-challenge/server"
	"github.com/calvinbui/teleport-sre-challenge/services"
	"github.com/calvinbui/teleport-sre-challenge/services/logger"
)

func main() {
	// initalise our logger first
	logger.Init()

	ctx := context.Background()

	// load configuration values from args and env vars
	conf, err := config.New()
	if err != nil {
		logger.Fatal("Error getting config", err)
	}

	// create clients using config
	svc, err := services.New(&conf)
	if err != nil {
		logger.Fatal("Error creating services", err)
	}

	// create the 'server' instance which runs our commands
	srv := server.New(&svc, &conf)

	// http health check
	go func() {
		if err := srv.StartHTTPHealthCheck(ctx); err != nil {
			logger.Fatal("Error serving HTTP health check", err)
		}
	}()

	// deployment replica count reconciliation
	go func() {
		if err := srv.WatchAndReconcileDeployments(ctx); err != nil {
			logger.Fatal("Error reconciling deployments", err)
		}
	}()

	// grpc api
	logger.Info(fmt.Sprintf("Server is listening on %v", conf.GRPCAddress))
	if err := srv.ServeGRPC(); err != nil {
		logger.Fatal("Error serving gRPC", err)
	}
}
