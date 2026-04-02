package app

import (
	"log/slog"

	"time"

	grpcapp "github.com/voznikaetnepriyazn/Autorization_service/internal/app/grpc"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func InitApp(log *slog.Logger, grpcPort int, tokenTTL time.Duration) *App {
	grpcApp := grpcapp.InitApp(log, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
