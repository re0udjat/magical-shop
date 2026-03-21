package data

import "time"

type Item struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Rarity    string    `json:"rarity"`
	Price     Currency  `json:"price"`
	CreatedAt time.Time `json:"-"`
}
