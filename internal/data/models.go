package data

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// All models will be wrapped inside the Models struct
type Models struct {
	Items ItemModel
	Users UserModel
}

// Returns a Model struct containing the initialized ItemModel
func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Items: ItemModel{DB: db},
		Users: UserModel{DB: db},
	}
}
