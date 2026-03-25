package data

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

// All models will be wrapped inside the Models struct
type Models struct {
	Items ItemModel
}

// Returns a Model struct containing the initialized ItemModel
func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Items: ItemModel{DB: db},
	}
}
