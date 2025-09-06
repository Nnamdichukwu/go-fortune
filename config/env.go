package config

import(
	"github.com/joho/godotenv"
)

func LoadEnvVars() error {
	_ = godotenv.Load()
	if err := LoadPostgresConfig(); err != nil{
		return err
	}
	return nil

}