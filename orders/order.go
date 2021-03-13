package orders

import (
	"time"
)

type Order struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Temp       string  `json:"temp"`
	ShelfLife  float64 `json:"shelfLife"`
	DecayRate  float64 `json:"decayRate"`
	CreateTime time.Time
}




