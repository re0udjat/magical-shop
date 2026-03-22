package data

import (
	"time"

	"github.com/re0udjat/magic-shop/internal/validator"
)

type Item struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Rarity    string    `json:"rarity"`
	Price     Currency  `json:"price"`
	CreatedAt time.Time `json:"-"`
}

func ValidateItem(v *validator.Validator, item *Item) {
	v.Check(item.Name != "", "name", "must be provided")
	v.Check(len(item.Name) <= 500, "name", "must not be more than 500 bytes long")

	v.Check(validator.PermittedValue(item.Rarity, "common", "uncommon", "rare", "mythic", "legendary"), "rarity", "must be a valid rarity")

	v.Check(item.Price > 0, "price", "must be a positive value")
	v.Check(item.Price <= 1_000_000_000, "price", "must not be more than 1,000,000,000 coins")

}
