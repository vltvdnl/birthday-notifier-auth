package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sso/internal/domain/models"
	"sso/pkg/errs"
	"time"
)

type Authorizer interface {
	NewToken(user models.User, app models.App, duration time.Duration) (string, error)
	VerifyPassword(hash []byte, passhash []byte) error
	GenerateHash(password []byte) ([]byte, error)
}
type UserSaver interface {
	SaveUser(ctx context.Context, email string, passhash []byte) (uid int64, err error)
}
type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
}
type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}
type Auth struct {
	log          *slog.Logger
	authorizer   Authorizer
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

func New(log *slog.Logger, authorizer Authorizer, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{log: log, authorizer: authorizer, userSaver: userSaver, userProvider: userProvider, appProvider: appProvider, tokenTTL: tokenTTL}
}

func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	const log_op = "Auth.Login"

	log := a.log.With(slog.String("op", log_op), slog.String("username", email))
	log.Info("attempting to login user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			log.Warn("user not found", err)
			return "", fmt.Errorf("%s: %w", log_op, errs.ErrUserNotFound)
		}
		log.Error("failed to get user", err)
		return "", fmt.Errorf("%s: %w", log_op, err)
	}

	if err := a.authorizer.VerifyPassword([]byte(user.Passhash), []byte(password)); err != nil {
		log.Warn("invalid credentials", err)
		return "", fmt.Errorf("%s: %w", log_op, errs.ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		log.Warn("invalid credential: no app with this id", err)
		return "", fmt.Errorf("%s: %w", log_op, err)
	}
	log.Info("user successfully logged")

	token, err := a.authorizer.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("Failed to generate token", err)
		return "", fmt.Errorf("%s: %w", log_op, err)
	}
	return token, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, pass string) (int64, error) {
	const log_op = "Auth.RegisterNewUser"

	log := a.log.With(slog.String("op", log_op), slog.String("email", email))
	log.Info("registering user")

	passHash, err := a.authorizer.GenerateHash([]byte(pass))
	if err != nil {
		log.Error("failed to generate password hash", err)
		return 0, fmt.Errorf("%s: %w", log_op, err)
	}
	if id, err := a.userSaver.SaveUser(ctx, email, passHash); err != nil {
		log.Error("failed to save user")
		return 0, fmt.Errorf("%s: %w", log_op, err)
	} else {
		return id, nil
	}
}
