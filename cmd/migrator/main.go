package main

import (
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // подключаем драйвер для PostgreSQL
	_ "github.com/golang-migrate/migrate/v4/source/file"       // поддержка чтения миграций из файлов
)

func main() {
	var user, password, host, port, dbname, migrationsPath string //TODO: too unsafe to use flags here; need some fixes
	var fl string
	flag.StringVar(&user, "user", "postgres", "pg user")
	flag.StringVar(&password, "password", "123", "pg password")
	flag.StringVar(&host, "host", "localhost", "pg host")
	flag.StringVar(&port, "port", "5432", "pg port")
	flag.StringVar(&dbname, "dbname", "postgres", "pg dbname")
	flag.StringVar(&migrationsPath, "mgpath", "./migrations", "migrations path")
	flag.StringVar(&fl, "fl", "up", "up or down")

	flag.Parse()
	// fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
	m, err := migrate.New("file://"+migrationsPath, fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname))
	if err != nil {
		panic("1." + err.Error())
	}
	if fl == "up" {
		if err := m.Up(); err != nil {
			panic(err)
		}
	} else {
		if err := m.Down(); err != nil {
			panic(err)
		}
	}
	fmt.Println("migrations successfully applied")

}
