package app

import (
	"log/slog"
	"sso/internal/app/grpc"
	"sso/internal/services/auth"
	"sso/pkg/jwt"
	"sso/pkg/storage"
	"time"
)

type App struct {
	GRPCServer *grpc.App
}

func New(log *slog.Logger, grpcPort int, pg_url string, tokenTTL time.Duration) *App {
	storage, err := storage.New(pg_url)
	if err != nil {
		panic(err)
	}
	jwtService := jwt.New()
	authService := auth.New(log, jwtService, storage, storage, storage, tokenTTL)
	grpcApp := grpc.New(log, authService, grpcPort)
	return &App{grpcApp}

}
