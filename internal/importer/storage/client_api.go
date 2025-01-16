package storage

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/client"
	"github.com/green-ecolution/tbz-csv-import-plugin/internal/entities"
)

type GreenEcolutionRepo struct {
	client *client.APIClient
}

func NewGreenEcolutionRepo(cfg *client.Configuration) *GreenEcolutionRepo {
	return &GreenEcolutionRepo{
		client: client.NewAPIClient(cfg),
	}
}

func (r *GreenEcolutionRepo) GetInfo(ctx context.Context) (*client.AppInfo, error) {
	info, _, err := r.client.InfoAPI.GetAppInfo(ctx).Execute()
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (r *GreenEcolutionRepo) GetTrees(ctx context.Context) ([]client.Tree, error) {
	trees, _, err := r.client.TreeAPI.GetAllTrees(ctx).Execute()
	if err != nil {
		return nil, err
	}
	return trees.Data, nil
}

func (r *GreenEcolutionRepo) CreateTrees(ctx context.Context, trees []*entities.Tree) error {
	for _, tree := range trees {
		_, _, err := r.client.TreeAPI.CreateTree(ctx).Body(client.TreeCreate{
			Description:  "Dieser Baum wurde von einem CSV-Import erstellt.",
			Latitude:     float32(tree.Latitude),
			Longitude:    float32(tree.Longitude),
			PlantingYear: tree.PlantingYear,
			Species:      tree.Species,
			TreeNumber:   tree.Number,
		}).Execute()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *GreenEcolutionRepo) UpdateTrees(ctx context.Context, trees []*entities.Tree) error {
	for _, tree := range trees {
		_, _, err := r.client.TreeAPI.UpdateTree(ctx, string(tree.TreeID)).Body(client.TreeUpdate{
			Description:  "Dieser Baum wurde von einem CSV-Import aktualisiert.",
			Latitude:     float32(tree.Latitude),
			Longitude:    float32(tree.Longitude),
			PlantingYear: tree.PlantingYear,
			Species:      tree.Species,
			TreeNumber:   tree.Number,
		}).Execute()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *GreenEcolutionRepo) DeleteTrees(ctx context.Context, treeIDs []entities.TreeID) error {
	for _, treeID := range treeIDs {
		_, err := r.client.TreeAPI.DeleteTree(ctx, string(treeID)).Execute()
		if err != nil {
			return err
		}
	}
	return nil
}
