package storage

import (
	"context"

	"database/sql"
	"sso/internal/domain/models"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func New(url string) (*Postgres, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &Postgres{db: db}, nil
}
func (r *Postgres) Stop() error {
	return r.db.Close()
}
func (r *Postgres) SaveUser(ctx context.Context, email string, passhash []byte) (uid int64, err error) {
	// TODO: validate if user is exists
	// const log_op = "Postgres.SaveUser"
	sqlstatement := `INSERT INTO users(email, pass_hash) VALUES ($1, $2)`
	_, err = r.db.ExecContext(ctx, sqlstatement, email, passhash)
	if err != nil {
		return 0, err
	}
	err = r.db.QueryRowContext(ctx, `SELECT id FROM users WHERE email = $1`, email).Scan(&uid)
	if err != nil {
		return 0, err
	}
	return uid, nil
}
func (r *Postgres) User(ctx context.Context, email string) (models.User, error) {
	// const log_op = "Postgres.User"

	sqlstatement := `SELECT id, email, pass_hash FROM users WHERE email = $1`
	var user models.User
	err := r.db.QueryRowContext(ctx, sqlstatement, email).Scan(&user.ID, &user.Email, &user.Passhash)
	if err != nil {
		return models.User{}, err // TODO if user is not registered
	}
	return user, nil
}
func (r *Postgres) App(ctx context.Context, id int) (models.App, error) {
	// const log_op = "Postgres.App"

	sqlstatement := `SELECT id, name, secret FROM apps WHERE id = $1`
	var app models.App

	err := r.db.QueryRow(sqlstatement, id).Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		return models.App{}, err //TODO if app is not found
	}
	return app, nil
}
