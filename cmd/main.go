package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
	"os/signal"
	"syscall"
	"urleater/internal/config"
	"urleater/internal/handlers"
	"urleater/internal/repository/postgresDB"
	"urleater/internal/service"
)

const port = ":8080"

func main() {
	serverCtx, serverCancel := context.WithCancel(context.Background())
	postgresConfig := config.Config{
		DB: config.DBConfig{
			PostgresHost:     "localhost",
			PostgresPort:     "5432",
			PostgresUser:     "postgres",
			PostgresPassword: "postgres",
			PostgresDatabase: "postgres",
		}}
	postgresPool := providePool(serverCtx, postgresConfig.PostgresURL(), true)

	// storage layer
	postgresStorage := postgresDB.NewStorage(postgresPool)

	// service layer
	srv := service.New(postgresStorage)

	// handlers layer
	e := handlers.GetRoutes(&handlers.Handlers{Service: srv})

	if err := e.Start(port); err != nil {
		log.Fatalf("could not start server: %v", err)
	}

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	serverCancel()
}

func providePool(ctx context.Context, url string, lazy bool) *pgxpool.Pool {
	poolConfig, err := pgxpool.ParseConfig(url)

	if err != nil {
		log.Fatal("Unable to parse DB config because " + err.Error())
	}

	poolConfig.LazyConnect = lazy

	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		log.Fatal("Unable to establish connection to " + url + " because " + err.Error())
	}

	return pool
}
