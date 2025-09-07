package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-version"

	"github.com/Nnamdichukwu/go-fortune/models"
)

func GetVersionById(ctx context.Context, db *sql.DB, id int) (*models.DbResponse, error) {
	query := `SELECT id, owner, repo, version  FROM packages WHERE id = $1`
	row := db.QueryRowContext(ctx, query, id)

	var resp models.DbResponse

	if err := row.Scan(&resp.ID, &resp.Owner, &resp.Repo, &resp.Version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &resp, nil

}

func GetVersionByOwner(ctx context.Context, db *sql.DB, owner string) (*models.DbResponse, error) {
	query := `SELECT id, owner, repo, version  FROM packages WHERE owner = $1`
	row := db.QueryRowContext(ctx, query, owner)

	var resp models.DbResponse

	if err := row.Scan(&resp.ID, &resp.Owner, &resp.Repo, &resp.Version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &resp, nil
}

func GetVersionByRepo(ctx context.Context, db *sql.DB, repo string) (*models.DbResponse, error) {
	query := `SELECT id, owner, repo, version  FROM packages WHERE repo = $1`
	row := db.QueryRowContext(ctx, query, repo)

	var resp models.DbResponse

	if err := row.Scan(&resp.ID, &resp.Owner, &resp.Repo, &resp.Version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &resp, nil
}

func GetVersionByOwnerAndRepo(ctx context.Context, db *sql.DB, owner string, repo string) (*models.DbResponse, error) {
	query := `SELECT id, owner, repo, version  FROM packages WHERE owner = $1 AND repo = $2`
	row := db.QueryRowContext(ctx, query, owner)

	var resp models.DbResponse

	if err := row.Scan(&resp.ID, &resp.Owner, &resp.Repo, &resp.Version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &resp, nil
}

func UpdateVersion(ctx context.Context, db *sql.DB, owner string, repo string, new_version string) error {
	resp, err := GetVersionByOwnerAndRepo(ctx, db, owner, repo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return err
	}

	current_version, err := version.NewVersion(resp.Version)
	if err != nil {
		return err
	}
	latest_version, err := version.NewVersion(new_version)
	if err != nil {
		return err
	}

	if !current_version.LessThan(latest_version) {
		return errors.New("the current version is more recent than the latest version")
	}
	query := `UPDATE packages SET version = $1, updated_at= $2 WHERE owner = $3 and repo = $4`
	res, err := db.ExecContext(ctx, query, latest_version, time.Now(), owner, repo)
	if err != nil {
		return fmt.Errorf("updated failed due to: %w", err)
	}
	rowNo, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rowNo == 0 {
		return errors.New("concurrent update occured")
	}

	return nil
}
