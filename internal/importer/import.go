package importer

import (
	"context"
	"slices"

	"github.com/green-ecolution/tbz-csv-import-plugin/internal/entities"
	"github.com/green-ecolution/tbz-csv-import-plugin/internal/importer/storage"
)

type ImportService struct {
	importRepo *storage.ImportRepositoryDB
	clientRepo *storage.GreenEcolutionRepo
}

func NewImportService() *ImportService {
	return &ImportService{}
}

func (i *ImportService) Import(ctx context.Context, trees []*entities.TreeImport) error {
	deleteQueue := make([]entities.TreeID, len(trees))
	createQueue := make([]*entities.TreeImport, len(trees))
	updateQueue := make([]*entities.TreeImport, len(trees))

	allImportedTrees, err := i.importRepo.GetAllTrees(ctx)
	if err != nil {
		return err
	}

	for _, csvTree := range trees {
		idx := slices.IndexFunc(allImportedTrees, func(tree entities.TreeImport) bool {
			return tree.Latitude == csvTree.Latitude && tree.Longitude == csvTree.Longitude
		})
		if idx == -1 {
			createQueue = append(createQueue, csvTree)
			continue
		}

		existingTree := allImportedTrees[idx]
		if existingTree.PlantingYear == csvTree.PlantingYear {
			csvTree.TreeID = existingTree.TreeID
			updateQueue = append(updateQueue, csvTree)
		} else {
			deleteQueue = append(deleteQueue, existingTree.TreeID)
			createQueue = append(createQueue, csvTree)
		}
	}

	err = i.importRepo.WithTx(ctx, func(ctx context.Context, tx *storage.ImportRepositoryTx) error {
		if err := tx.CreateTrees(ctx, createQueue); err != nil {
			return err
		}

		if err := tx.UpdateTrees(ctx, updateQueue); err != nil {
			return err
		}

		if err := tx.DeleteTreesByID(ctx, deleteQueue); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
