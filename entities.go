package main

import (
	"fmt"
	"slices"

	"github.com/green-ecolution/green-ecolution-backend/pkg/client"
)

type Tree struct {
	client.Tree
	ObjectID int
}

type CsvTree struct {
	Area         string
	Street       string
	TreeNumber   string
	Species      string
	Hochwert     float64
	Rechtswert   float64
	PlantingYear int
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
