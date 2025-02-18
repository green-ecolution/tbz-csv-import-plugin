package main

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
)

type SyncTrees struct {
	csvTrees []CsvTree
	client   *GreenEcolutionClient
}

func NewSyncTrees(csvTrees []CsvTree, client *GreenEcolutionClient) *SyncTrees {
	return &SyncTrees{
		csvTrees: csvTrees,
		client:   client,
	}
}

func (s *SyncTrees) Sync(ctx context.Context) (TreeImportResponse, error) {
	slog.Info("sync tbz register trees to green ecolution backend")

	mapCsvTrees, err := TreesFromBatch(s.csvTrees)
	if err != nil {
		slog.Error("failed to map csv trees to internal trees", "error", err)
		return TreeImportResponse{}, err
	}

	geTrees, err := s.client.GetAll(ctx)
	if err != nil {
		slog.Error("failed to get trees from green ecolution backend", "error", err)
		return TreeImportResponse{}, err
	}

	slices.SortFunc(mapCsvTrees, func(a, b Tree) int {
		return a.ObjectID - b.ObjectID
	})

	slices.SortFunc(geTrees, func(a, b Tree) int {
		return a.ObjectID - b.ObjectID
	})

	createdQueue := make([]Tree, 0)
	updateQueue := make([]Tree, 0)
	archiveQueue := make([]Tree, 0)

	idxCsvTrees := 0
	idxGeTrees := 0

	for idxCsvTrees < len(mapCsvTrees) || idxGeTrees < len(geTrees) {
		if idxCsvTrees == len(mapCsvTrees) {
			archiveQueue = append(archiveQueue, geTrees[idxCsvTrees:]...)
			break
		}

		if idxGeTrees == len(geTrees) {
			createdQueue = append(createdQueue, mapCsvTrees[idxGeTrees:]...)
			break
		}

		csvTree := mapCsvTrees[idxCsvTrees]
		geTree := geTrees[idxGeTrees]

		fmt.Println(csvTree.ObjectID, geTree.ObjectID)

		if csvTree.ObjectID == geTree.ObjectID {
			if updatedTree, ok := s.checkDiff(csvTree, geTree); !ok {
				updateQueue = append(updateQueue, updatedTree)
			}
			idxGeTrees++
			idxCsvTrees++
			continue
		}

		if csvTree.ObjectID < geTree.ObjectID {
			createdQueue = append(createdQueue, csvTree)
			idxCsvTrees++
			continue
		}

		if csvTree.ObjectID > geTree.ObjectID {
			archiveQueue = append(archiveQueue, geTree)
			idxGeTrees++
			continue
		}
	}

	fmt.Println("create queue")
	fmt.Printf("%+v\n", createdQueue)

	fmt.Println("update queue")
	fmt.Printf("%+v\n", updateQueue)

	fmt.Println("archive queue")
	fmt.Printf("%+v\n", archiveQueue)

	for _, e := range createdQueue {
		if err := s.client.Create(ctx, e); err != nil {
			slog.Warn("failed to create tree in green ecolution backend", "error", err, "register_id", e.ObjectID)
		}
	}

	for _, e := range updateQueue {
		if err := s.client.Update(ctx, e.Id, e); err != nil {
			slog.Warn("failed to update tree in green ecolution backend", "error", err, "register_id", e.ObjectID, "tree_id", e.Id)
		}
	}

	for _, e := range archiveQueue {
		if err := s.client.Archive(ctx, e.Id); err != nil {
			slog.Warn("failed to archive tree in green ecolution backend", "error", err, "register_id", e.ObjectID, "tree_id", e.Id)
		}
	}

	importedTrees := make([]TreeImport, len(archiveQueue)+len(updateQueue)+len(createdQueue))
	importedTrees = append(importedTrees, Map(createdQueue, func(t Tree) TreeImport {
		return TreeImport{
			Tree:       t,
			ImportType: ImportTypeCreate,
		}
	})...)

	importedTrees = append(importedTrees, Map(updateQueue, func(t Tree) TreeImport {
		return TreeImport{
			Tree:       t,
			ImportType: ImportTypeUpdate,
		}
	})...)

	importedTrees = append(importedTrees, Map(archiveQueue, func(t Tree) TreeImport {
		return TreeImport{
			Tree:       t,
			ImportType: ImportTypeArchive,
		}
	})...)

	slices.SortFunc(importedTrees, func(a, b TreeImport) int {
		return a.Tree.ObjectID - b.Tree.ObjectID
	})

	return TreeImportResponse{
		ImportedTrees: importedTrees,
		Raw:           s.csvTrees,
	}, nil
}

func (s *SyncTrees) checkDiff(new, old Tree) (Tree, bool) {
	if new.Number == old.Number ||
		new.Latitude == old.Latitude ||
		new.Longitude == old.Longitude ||
		new.Species == old.Species ||
		new.PlantingYear == old.PlantingYear {
		return new, true
	} else {
		old.Number = new.Number
		old.Latitude = new.Latitude
		old.Longitude = new.Longitude
		old.Species = new.Species
		old.PlantingYear = new.PlantingYear

		return old, false
	}
}
