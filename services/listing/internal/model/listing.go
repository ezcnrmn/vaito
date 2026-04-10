package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
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

type Pagination struct {
	Page          int32
	Size          int32
	Sort          string
	SortDirection string
	Filter        struct {
		Status *string
		UserID *int64
	}
	Total int32
}

func (p Pagination) orderBy() string {
	return fmt.Sprintf("%s %s", p.Sort, p.SortDirection)
}

func (p Pagination) limit() int32 {
	if p.Size < 1 || p.Size > 100 {
		return 20
	}
	return p.Size
}

func (p Pagination) offset() int32 {
	if p.Page < 1 {
		return 0
	}
	return (p.Page - 1) * p.limit()
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
	SET title = $1, description = $2, price = $3, category_id = $4, version = version+1, status_id = (SELECT id FROM listing_statuses WHERE name='Draft'), published_at=NULL
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

func (lm ListingModel) GetListing(ctx context.Context, id int64, userID *int64, status *string) (*Listing, error) {
	where := []string{fmt.Sprintf("l.id=%d", id)}
	if userID != nil {
		where = append(where, fmt.Sprintf("l.user_id=%d", *userID))
	}
	if status != nil {
		where = append(where, fmt.Sprintf("l.status_id=(SELECT id FROM listing_statuses WHERE name='%s')", *status))
	}
	whereCond := strings.Join(where, " AND ")

	query := fmt.Sprintf(`
	SELECT l.id, l.title, l.description, l.category_id, c.name, l.user_id, s.name, l.price, l.created_at, l.published_at, l.version
	FROM listings l
	JOIN listing_statuses s ON l.status_id = s.id
	JOIN categories c ON l.category_id = c.id
	WHERE %s;`, whereCond)

	listing := &Listing{}
	err := lm.DB.QueryRowContext(ctx, query).Scan(
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

func (lm ListingModel) GetListings(ctx context.Context, pagination *Pagination) ([]Listing, error) {
	where := []string{"TRUE"}
	if pagination.Filter.Status != nil {
		where = append(where, fmt.Sprintf("s.name='%s'", *pagination.Filter.Status))
	}
	if pagination.Filter.UserID != nil {
		where = append(where, fmt.Sprintf("l.user_id=%d", *pagination.Filter.UserID))
	}
	whereCond := strings.Join(where, " AND ")

	query := fmt.Sprintf(`
	SELECT count(*) OVER(), l.id, l.title, l.description, l.category_id, c.name, l.user_id, s.name, l.price, l.created_at, l.published_at, l.version
	FROM listings l
	JOIN listing_statuses s ON l.status_id = s.id
	JOIN categories c ON l.category_id = c.id
	WHERE %s
	ORDER BY %s
	LIMIT $1 OFFSET $2`, whereCond, pagination.orderBy())

	args := []any{pagination.limit(), pagination.offset()}

	rows, err := lm.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var listings []Listing
	for rows.Next() {
		var l Listing
		err := rows.Scan(
			&pagination.Total,
			&l.ID,
			&l.Title,
			&l.Description,
			&l.Cetegory.ID,
			&l.Cetegory.Name,
			&l.UserID,
			&l.Status,
			&l.Price,
			&l.CreatedAt,
			&l.PublishedAt,
			&l.Version,
		)
		if err != nil {
			return nil, err
		}
		listings = append(listings, l)
	}

	return listings, nil
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
	SET
		status_id = (SELECT id FROM listing_statuses WHERE name = $1),
		version = version+1,
		published_at = (CASE WHEN $1 = 'Active' THEN $2 ELSE NULL END)::timestamptz
	WHERE id = $3 AND version = $4;`

	args := []any{status, time.Now(), listing.ID, listing.Version}

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
