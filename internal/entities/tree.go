package entities

import "time"

type Tree struct {
	TreeID       TreeID           `db:"id"`
	CreatedAt    time.Time        `db:"created_at"`
	UpdatedAt    time.Time        `db:"updated_at"`
	Area         TreeArea         `db:"area"`
	Number       TreeNumber       `db:"tree_number"`
	Species      TreeSpecies      `db:"species"`
	Latitude     TreeLatitude     `db:"latitude"`
	Longitude    TreeLongitude    `db:"longitude"`
	PlantingYear TreePlantingYear `db:"planting_year"`
	Street       TreeStreet       `db:"street"`
}

type TreeArea = string
type TreeNumber = string
type TreeSpecies = string
type TreeLatitude = float64
type TreeLongitude = float64
type TreePlantingYear = int32
type TreeStreet = string
type TreeID = int32

type Import struct {
	ID        ImportID  `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UserID    UserID    `db:"user_id"`
	RawCSV    RawCSV    `db:"raw_csv"`
}

type ImportID = int32
type UserID = string
type RawCSV = string

