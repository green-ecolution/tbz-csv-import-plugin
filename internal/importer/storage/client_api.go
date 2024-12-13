package storage

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/client"
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


