package sso

import (
	"log/slog"
	"os"

	"github.com/voznikaetnepriyazn/Autorization_service/internal/app"
	"github.com/voznikaetnepriyazn/Autorization_service/internal/config"
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

	application := app.InitApp(log, int(cfg.GRPCApp.Port), cfg.TokenTTL)

	application.GRPCSrv.MustRun()
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
