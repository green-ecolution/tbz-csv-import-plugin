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
		return Tree{
			ObjectID: (int(g.X) << 8) + int(g.Y),
			Tree: client.Tree{
				PlantingYear: int32(t.PlantingYear),
				Species:      t.Species,
				Number:       t.TreeNumber,
				Description:  fmt.Sprint("%s %s", t.Area, t.Street),
				Latitude:     float32(g.X),
				Longitude:    float32(g.Y),
			},
		}
	})

	return slices.Collect(treeSeq), nil
}
