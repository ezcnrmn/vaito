package model

import (
	"database/sql"
)

type Listing struct{}

type ListingModel struct {
	DB *sql.DB
}

func (um ListingModel) Insert(user *Listing) error {
	return nil
}
