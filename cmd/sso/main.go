package main

import (
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"sso/pkg/log"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log := log.New(cfg.Env)

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.GRPC.Timeout) //replace storage path

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	application.GRPCServer.Stop()
	log.Info("gracefully stop")
}
