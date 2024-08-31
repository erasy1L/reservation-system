package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"room-reservation/internal/handler"
	"room-reservation/internal/repository"
	"room-reservation/pkg/log"
	"room-reservation/pkg/server"
	"syscall"
	"time"
)

func main() {
	logger := log.LoggerFromContext(context.Background())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	fmt.Println("Press Ctrl+C to exit")

	connString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	reservationRepo, err := repository.NewReservationRepository(context.Background(), connString)
	if err != nil {
		logger.Fatal().Err(err).Msg("error intializing reservation repository")
	}

	reservationHTTPHandler := handler.NewReservationHandler(reservationRepo)

	httpServer := server.New(reservationHTTPHandler.HTTP, os.Getenv("APP_PORT"))

	fmt.Println("Swagger is accessible at http://localhost:" + os.Getenv("APP_PORT") + "/swagger/index.html")

	if err := httpServer.Start(); err != nil {
		logger.Fatal().Err(err).Msg("error starting http server")
	}

	<-stop
	fmt.Println("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	reservationRepo.Close()

	if err := httpServer.Stop(ctx); err != nil {
		logger.Fatal().Err(err).Msg("error stopping server")
	}

	fmt.Println("server successfully shutdown")
}
