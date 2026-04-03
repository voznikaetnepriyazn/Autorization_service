package main

//todo interseptor - like middleware, graceful shutdown for db ping

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/voznikaetnepriyazn/Autorization_service/internal/app"
	"github.com/voznikaetnepriyazn/Autorization_service/internal/config"
	"github.com/voznikaetnepriyazn/Autorization_service/internal/grpc/auth"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting aplication",
		slog.String("env", cfg.Env),
		slog.Any("cfg", cfg),
		slog.Int("port", int(cfg.GRPCApp.Port)),
	)

	var authService auth.Auth

	application := app.InitApp(log, authService, int(cfg.GRPCApp.Port), cfg.TokenTTL)

	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	signal := <-stop //blocking operation

	log.Info("stopping application", slog.String("signal", signal.String()))

	application.GRPCSrv.Stop()

	log.Info("application stop")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
