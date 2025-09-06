package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Nnamdichukwu/go-fortune/config"
	"github.com/Nnamdichukwu/go-fortune/models"
)
var PostgresDB *sql.DB

func ConnectPostgresDB(credentials config.Postgres) error{
	connStr := fmt.Sprintf("postgres://%s:%s@db:%s/%s?sslmode=disable",
	credentials.Host, credentials.Password, credentials.Port, credentials.Name)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		 return err
	}
	defer db.Close()

	if err = db.Ping(); err != nil{
		 return err

	}
	PostgresDB = db
	return nil
}

func InsertIntoPostgres(ctx context.Context, db *sql.DB,response models.Response, table string) (int,error) {
	query := fmt.Sprintf("INSERT INTO %s(owner, repo, version, created_at, updated_at)", table)
	var pk int
	err := db.QueryRowContext(ctx, query, response.Owner,response.Repo, response.Version, response.CreatedAt, response.UpdatedAt).Scan(&pk)
	if err != nil{

		return 0, err
	}
	return pk, nil

}

func CreatePackagesTable(ctx context.Context,db *sql.DB) error{
	
	query := `CREATE TABLE IF NOT EXISTS packages (
		id SERIAL PRIMARY KEY,
		owner VARCHAR(100) NOT NULL,
		repo VARCHAR(100) NOT NULL,
		version VARCHAR(100) NOT NULL
		created_at TIMESTAMP DEFAULT NOW()
		updated_at TIMESTAMP DEFAULT NOW()
	)`
	_, err := db.ExecContext(ctx,query)
	if err != nil {
		return err
	}
	return nil

}
