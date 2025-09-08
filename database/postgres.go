package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Nnamdichukwu/go-fortune/config"
	"github.com/Nnamdichukwu/go-fortune/models"
	_ "github.com/lib/pq"
)

var PostgresDB *sql.DB

func ConnectPostgresDB(credentials config.Postgres) error {
	connStr := fmt.Sprintf("postgres://%s:%s@db:%s/%s?sslmode=disable",
		credentials.Host, credentials.Password, credentials.Port, credentials.Name)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return errors.New("failed to connect to postgress")
	}
	

	if err = db.Ping(); err != nil {
		return errors.New("cannot ping db")

	}
	PostgresDB = db
	return nil
}

func InsertIntoPostgres(ctx context.Context, db *sql.DB, response models.Response) (int, error) {
	query := `INSERT INTO packages(owner, repo, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)
	RETURNING id`
	var pk int
	err := db.QueryRowContext(ctx, query, response.Owner, response.Repo, response.Version, response.CreatedAt, response.UpdatedAt).Scan(&pk)
	if err != nil {

		return 0, err
	}
	return pk, nil

}

func CreatePackagesTable(ctx context.Context, db *sql.DB) error {

	query := `CREATE TABLE IF NOT EXISTS packages (
		id SERIAL PRIMARY KEY,
		owner VARCHAR(100) NOT NULL,
		repo VARCHAR(100) NOT NULL,
		version VARCHAR(100) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	)`
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil

}
