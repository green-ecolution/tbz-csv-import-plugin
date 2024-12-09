package entities

type TreeImport struct {
	Area         string `validate:"required"`
	Number       string `validate:"required"`
	Species      string
	Latitude     float64 `validate:"required,max=90,min=-90"`
	Longitude    float64 `validate:"required,max=180,min=-180"`
	PlantingYear int32   `validate:"gt=0"`
	Street       string  `validate:"required"`
	TreeID       int32
}

