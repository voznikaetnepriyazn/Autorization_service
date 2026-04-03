package app

import (
	"log/slog"

	"time"

	grpcapp "github.com/voznikaetnepriyazn/Autorization_service/internal/app/grpc"
	"github.com/voznikaetnepriyazn/Autorization_service/internal/grpc/auth"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func InitApp(log *slog.Logger, authService auth.Auth, grpcPort int, tokenTTL time.Duration) *App {
	grpcApp := grpcapp.InitApp(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
