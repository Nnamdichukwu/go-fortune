package config

import (
	"errors"
	"os"
)

type Postgres struct {
	Host     string
	Port     string
	Password string
	Name     string
}

var PostgresConfig Postgres

func LoadPostgresConfig() error {
	host, exist := os.LookupEnv("DB_HOST")
	if !exist {
		return errors.New("DB_HOST is not set in .env")
	}
	port, exist := os.LookupEnv("DB_PORT")
	if !exist {
		return errors.New("PORT is not set in .env")
	}

	pwd, exist := os.LookupEnv("DB_PASSWORD")
	if !exist {
		return errors.New("DB_PASSWORD is not set in .env")
	}
	name, exist := os.LookupEnv("DB_NAME")
	if !exist {
		return errors.New("DB_NAME is not set in .env")
	}
	PostgresConfig = Postgres{
		Host:     host,
		Port:     port,
		Password: pwd,
		Name:     name,
	}
	return nil

}
