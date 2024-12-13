package entities

type TreeImport struct {
	TreeID       TreeID           `db:"id"`
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
