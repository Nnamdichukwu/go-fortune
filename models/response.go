package models

import (
	"time"
)

type Response struct {
	Owner     string    `json:"owner"`
	Repo      string    `json:"repo"`
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
type PackageResponse struct {
	ID int
	Response
}

type VersionUpdate struct {
	Owner string `json:"owner"`
	Repo string `json:"repo"`
	Version string `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
}