package data

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/re0udjat/magic-shop/internal/validator"
)

type Item struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Rarity    string    `json:"rarity"`
	Price     Currency  `json:"price"`
	CreatedAt time.Time `json:"-"`
	Version   int       `json:"-"`
}

func ValidateItem(v *validator.Validator, item *Item) {
	v.Check(item.Name != "", "name", "must be provided")
	v.Check(len(item.Name) <= 500, "name", "must not be more than 500 bytes long")

	v.Check(validator.PermittedValue(item.Rarity, "common", "uncommon", "rare", "mythic", "legendary"), "rarity", "must be a valid rarity")

	v.Check(item.Price > 0, "price", "must be a positive value")
	v.Check(item.Price <= 1_000_000_000, "price", "must not be more than 1,000,000,000 coins")

}

// Model which wraps a pgxpool.Pool connection pool
type ItemModel struct {
	DB *pgxpool.Pool
}

func (m ItemModel) Insert(item *Item) error {
	query := `
		INSERT INTO items (name, rarity, price)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	args := []any{item.Name, item.Rarity, item.Price}

	// Create a context with a timeout of 3 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRow(ctx, query, args...).Scan(&item.ID, &item.CreatedAt)
}

func (m ItemModel) Get(id int64) (*Item, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, name, rarity, price, created_at, version
		FROM items
		WHERE id = $1`

	var item Item

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, query, id).Scan(&item.ID, &item.Name, &item.Rarity, &item.Price, &item.CreatedAt, &item.Version)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &item, nil
}

func (m ItemModel) Update(item *Item) error {
	query := `
		UPDATE items
		SET name = $1, rarity = $2, price = $3, version = version + 1
		WHERE id = $4 AND version = $5
		RETURNING version`

	args := []any{item.Name, item.Rarity, item.Price, item.ID, item.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, query, args...).Scan(&item.Version)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m ItemModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM items WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
