package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Nnamdichukwu/go-fortune/models"
	"github.com/Nnamdichukwu/go-fortune/requests"
	"github.com/hashicorp/go-version"
)

func GetVersionById(ctx context.Context, db *sql.DB, id int) (*models.PackageResponse, error) {
	query := `SELECT id, owner, repo, version  FROM packages WHERE id = $1`
	row := db.QueryRowContext(ctx, query, id)

	var resp models.PackageResponse

	if err := row.Scan(&resp.ID, &resp.Owner, &resp.Repo, &resp.Version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &resp, nil

}

func GetVersionByOwner(ctx context.Context, db *sql.DB, owner string) (*models.PackageResponse, error) {
	query := `SELECT id, owner, repo, version  FROM packages WHERE owner = $1`
	row := db.QueryRowContext(ctx, query, owner)

	var resp models.PackageResponse

	if err := row.Scan(&resp.ID, &resp.Owner, &resp.Repo, &resp.Version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &resp, nil
}

func GetVersionByRepo(ctx context.Context, db *sql.DB, repo string) (*models.PackageResponse, error) {
	query := `SELECT id, owner, repo, version  FROM packages WHERE repo = $1`
	row := db.QueryRowContext(ctx, query, repo)

	var resp models.PackageResponse

	if err := row.Scan(&resp.ID, &resp.Owner, &resp.Repo, &resp.Version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &resp, nil
}

func GetVersionByOwnerAndRepo(ctx context.Context, db *sql.DB, r requests.Request) (*models.PackageResponse, error) {
	query := `SELECT id, owner, repo, version  FROM packages WHERE owner = $1 AND repo = $2`
	row := db.QueryRowContext(ctx, query, r.Owner, r.Repo)

	var resp models.PackageResponse

	if err := row.Scan(&resp.ID, &resp.Owner, &resp.Repo, &resp.Version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &resp, nil
}

func UpdateVersion(ctx context.Context, db *sql.DB, current_version string, ver models.VersionUpdate) error {
	
	existing_version, err := version.NewVersion(current_version)
	if err != nil {
		return err
	}
	latest_version, err := version.NewVersion(ver.Version)
	if err != nil {
		return err
	}

	if !existing_version.LessThan(latest_version) {
		return nil
	}
	query := `UPDATE packages SET version = $1, updated_at= $2 WHERE owner = $3 and repo = $4`
	res, err := db.ExecContext(ctx, query, latest_version.String(),ver.UpdatedAt, ver.Owner, ver.Repo)
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
