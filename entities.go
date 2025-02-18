package main

import (
	"fmt"
	"slices"

	"github.com/green-ecolution/green-ecolution-backend/pkg/client"
)

type Tree struct {
	client.Tree
	ObjectID int `json:"objectId"`
}

type CsvTree struct {
	Area         string  `json:"area"`
	Street       string  `json:"street"`
	TreeNumber   string  `json:"treeNumber"`
	Species      string  `json:"species"`
	Hochwert     float64 `json:"hochwert"`
	Rechtswert   float64 `json:"rechtswert"`
	PlantingYear int     `json:"plantingYear"`
}

type ImportType string

const (
	ImportTypeCreate  ImportType = "create"
	ImportTypeUpdate             = "update"
	ImportTypeArchive            = "archive"
)

type TreeImport struct {
	Tree       Tree       `json:"tree"`
	ImportType ImportType `json:"importType"`
}

type TreeImportResponse struct {
	ImportedTrees []TreeImport `json:"importedTrees"`
	Raw           []CsvTree    `json:"rawTrees"`
}

func TreesFromBatch(csvTrees []CsvTree) ([]Tree, error) {
	transformer, err := NewGeoTransformer(cfg.CsvUsedEpsg, cfg.CsvToEpsg)
	if err != nil {
		return nil, err
	}
	defer transformer.Destroy()

	geoPointsSeq := MapIter(slices.Values(csvTrees), func(t CsvTree) GeoPoint {
		return GeoPoint{
			X: t.Hochwert,
			Y: t.Rechtswert,
		}
	})

	geoPoints := slices.Collect(geoPointsSeq)
	geoPoints, err = transformer.TransformBatch(geoPoints)
	if err != nil {
		return nil, err
	}

	treeSeq := MapIter21(Zip(slices.Values(csvTrees), slices.Values(geoPoints)), func(t CsvTree, g GeoPoint) Tree {
		lat, long := float32(g.X), float32(g.Y)
		return Tree{
			ObjectID: CombineToID(lat, long),
			Tree: client.Tree{
				PlantingYear: int32(t.PlantingYear),
				Species:      t.Species,
				Number:       t.TreeNumber,
				Description:  fmt.Sprintf("%s %s", t.Area, t.Street),
				Latitude:     lat,
				Longitude:    long,
			},
		}
	})

	return slices.Collect(treeSeq), nil
}
