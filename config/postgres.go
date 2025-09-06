package config

import (
	"errors"
	"os"
)

type Mysql struct {
	Host string
	Port string
	User string
	Password string
	Name string
}
var MysqlConfig Mysql 

func LoadPostgresConfig() error {
	host, exist := os.LookupEnv("PG_HOST")
	if !exist{
		return errors.New("PG_HOST is not set in .env")
	}
	port, exist := os.LookupEnv("DB_PORT")
	if !exist{
		return errors.New("PORT is not set in .env")
	}
	user, exist := os.LookupEnv("DB_USER")
	if !exist{
		return errors.New("DB_USER is not set in .env")
	}
	pwd, exist := os.LookupEnv("DB_PASSWORD")
	if !exist{
		return errors.New("DB_PASSWORD is not set in .env")
	}
	name, exist := os.LookupEnv("DB_NAME")
	if !exist{
		return errors.New("PG_HOST is not set in .env")
	}
	MysqlConfig = Mysql{
		Host: host, 
		Port: port, 
		User: user, 
		Password: pwd, 
		Name: name,

	}
	return nil







}