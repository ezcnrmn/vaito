package model

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrListingNotFound = errors.New("listing not found")
	ErrEditConflict    = errors.New("edit conflict")
)

type Listing struct {
	ID          int64
	Title       string
	Description string
	Cetegory    Category
	UserID      int64
	Status      string
	Price       int64
	CreatedAt   time.Time
	PublishedAt *time.Time
	Version     int
}

type Category struct {
	ID   int64
	Name string
}

type ListingModel struct {
	DB *sql.DB
}

func (lm ListingModel) Insert(ctx context.Context, listing *Listing) error {
	query := `
	INSERT INTO listings (title, description, category_id, user_id, price, created_at, status_id)
	VALUES ($1, $2, $3, $4, $5, $6, (SELECT id FROM listing_statuses WHERE name='Draft'))
	RETURNING id, (SELECT name FROM categories WHERE id=category_id), (SELECT name FROM listing_statuses WHERE id=status_id), created_at`

	args := []any{listing.Title, listing.Description, listing.Cetegory.ID, listing.UserID, listing.Price, time.Now()}

	err := lm.DB.QueryRowContext(ctx, query, args...).Scan(&listing.ID, &listing.Cetegory.Name, &listing.Status, &listing.CreatedAt)

	// TODO: consider check constraints

	return err
}

func (lm ListingModel) UpdateListing(ctx context.Context, listing *Listing) error {
	query := `
	UPDATE listings
	SET title = $1, description = $2, price = $3, category_id = $4, version = version+1, status_id = (SELECT id FROM listing_statuses WHERE name='Draft')
	WHERE id = $5 AND version = $6
	RETURNING version, (SELECT name FROM listing_statuses WHERE id=status_id);`

	args := []any{listing.Title, listing.Description, listing.Price, listing.Cetegory.ID, listing.ID, listing.Version}

	err := lm.DB.QueryRowContext(ctx, query, args...).Scan(&listing.Version, &listing.Status)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrEditConflict
		default:
			return err
		}
	}

	return err
}

func (lm ListingModel) GetListing(ctx context.Context, id int64) (*Listing, error) {
	query := `
	SELECT l.id, l.title, l.description, l.category_id, c.name, l.user_id, s.name, l.price, l.created_at, l.published_at, l.version
	FROM listings l
	JOIN listing_statuses s ON l.status_id = s.id
	JOIN categories c ON l.category_id = c.id
	WHERE l.id = $1;`

	args := []any{id}

	listing := &Listing{}
	err := lm.DB.QueryRowContext(ctx, query, args...).Scan(
		&listing.ID,
		&listing.Title,
		&listing.Description,
		&listing.Cetegory.ID,
		&listing.Cetegory.Name,
		&listing.UserID,
		&listing.Status,
		&listing.Price,
		&listing.CreatedAt,
		&listing.PublishedAt,
		&listing.Version,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrListingNotFound
		default:
			return nil, err
		}
	}

	return listing, nil
}

func (lm ListingModel) DeleteListing(ctx context.Context, id, userID int64) error {
	query := `
	DELETE FROM listings
	WHERE id = $1 AND user_id = $2;`

	args := []any{id, userID}

	result, err := lm.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrListingNotFound
	}

	return nil
}

func (lm ListingModel) UpdateListingStatus(ctx context.Context, listing *Listing, status string) error {
	query := `
	UPDATE listings
	SET status_id = (SELECT id FROM listing_statuses WHERE name = $1), version = version+1
	WHERE id = $2 AND version = $3;`

	args := []any{status, listing.ID, listing.Version}

	result, err := lm.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrListingNotFound
	}

	return nil
}

func (lm ListingModel) GetCategories(ctx context.Context) (*[]Category, error) {
	query := `
	SELECT id, name
	FROM categories;`

	rows, err := lm.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var categories []Category
	for rows.Next() {
		var c Category
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return &categories, nil
}
