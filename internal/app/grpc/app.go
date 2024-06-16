package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sso/pkg/grpc/auth"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, authService auth.Auth, port int) *App {
	loggerOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}
	recoverOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))
			return status.Errorf(codes.Internal, "internal error")
		}),
	}
	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoverOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(log), loggerOpts...),
	))
	auth.Register(gRPCServer, authService)
	return &App{log: log, gRPCServer: gRPCServer, port: port}

}
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}
func (a *App) Run() error {
	const log_op = "App.Run"
	con, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", log_op, err)
	}
	a.log.Info("grpc server started", slog.String("addr", con.Addr().String()))
	if err := a.gRPCServer.Serve(con); err != nil {
		return fmt.Errorf("%s: %w", log_op, err)
	}
	return nil
}

func (a *App) Stop() {
	const log_op = "App.Stop"

	a.log.With(slog.String("op", log_op)).Info("stopping grpc server", slog.Int("port", a.port))
	a.gRPCServer.GracefulStop()
}
