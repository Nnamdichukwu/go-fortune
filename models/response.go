package models

import (
	"time"

)

type Request struct{
	Owner    string   `json:"owner"`
	Repo	 string   `json:"repo"`
	Version  string   `json:"version"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

}