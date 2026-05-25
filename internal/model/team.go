package model

type Team struct {
	ID       int     `json:"id" db:"id"`
	Name     string  `json:"name" db:"name"`
	Strength float32 `json:"strength" db:"strength"`
}
