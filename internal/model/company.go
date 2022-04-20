package model

import "time"

type Company struct {
	ID        int64      `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Code      string     `json:"code" db:"code"`
	Country   string     `json:"country" db:"country"`
	Website   string     `json:"website" db:"website"`
	Phone     string     `json:"phone" db:"phone"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}
