package models

import "time"

type Stats struct {
	ID        uint `json:"id" gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Storage      int `json:"storage"`
	CurrentOrder int `json:"currentOrder"`
	Debt         int `json:"debt"`
	GamerOrder   int `json:"gamerOrder"`
	Penalty      float32 `json:"penalty"`
	Step         int `json:"step"`
	GamerID      uint `gorm:"index"`
}
