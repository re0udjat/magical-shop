package data

import "time"

type Item struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Rarity    string    `json:"rarity"`
	Price     int64     `json:"price"`
	CreatedAt time.Time `json:"-"`
}
