package model

import (
	"time"
)

type User struct{
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	DisplayName string    `db:"display_name"`
	Bio         string    `db:"bio"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
