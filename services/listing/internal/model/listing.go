package model

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrListingNotFound = errors.New("listing not found")

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
}

type Category struct {
	ID   int64
	Name string
}

type ListingModel struct {
	DB *sql.DB
}

func (lm ListingModel) Insert(listing *Listing) error {
	query := `
	INSERT INTO listings (title, description, category_id, user_id, price, created_at, status_id)
	VALUES ($1, $2, $3, $4, $5, $6, (SELECT id FROM listing_statuses WHERE name='Draft'))
	RETURNING id, (SELECT name FROM categories WHERE id=category_id), (SELECT name FROM listing_statuses WHERE id=status_id), created_at`

	args := []any{listing.Title, listing.Description, listing.Cetegory.ID, listing.UserID, listing.Price, time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := lm.DB.QueryRowContext(ctx, query, args...).Scan(&listing.ID, &listing.Cetegory.Name, &listing.Status, &listing.CreatedAt)

	// TODO: check constraints

	return err
}

func (lm ListingModel) GetListing(id int64) (*Listing, error) {
	query := `
	SELECT l.id, l.title, l.description, l.category_id, c.name, l.user_id, s.name, l.price, l.created_at, l.published_at
	FROM listings l
	JOIN listing_statuses s ON l.status_id = s.id
	JOIN categories c ON l.category_id = c.id
	WHERE l.id = $1;`

	args := []any{id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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

func (lm ListingModel) GetCategories() (*[]Category, error) {
	query := `
	SELECT id, name
	FROM categories;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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
