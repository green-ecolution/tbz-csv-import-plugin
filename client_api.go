package main

import (
	"context"
	"fmt"
	"slices"

	"github.com/green-ecolution/green-ecolution-backend/pkg/client"
)

type GreenEcolutionClient struct {
	client   *client.APIClient
	provider string
}

func NewGreenEcolutionRepo(cfg *client.Configuration, provider string) *GreenEcolutionClient {
	return &GreenEcolutionClient{
		client:   client.NewAPIClient(cfg),
		provider: provider,
	}
}

func (r *GreenEcolutionClient) GetInfo(ctx context.Context) (*client.AppInfo, error) {
	info, _, err := r.client.InfoAPI.GetAppInfo(ctx).Execute()
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (r *GreenEcolutionClient) GetAll(ctx context.Context) ([]Tree, error) {
	trees, _, err := r.client.TreeAPI.GetAllTrees(ctx).Provider(r.provider).Execute()
	if err != nil {
		return nil, err
	}

	mapSeq := MapIter(slices.Values(trees.Data), func(t client.Tree) Tree {
		return Tree{
			Tree:     t,
			ObjectID: CombineToID(t.Latitude, t.Longitude),
		}
	})

	return slices.Collect(mapSeq), nil
}

func (r *GreenEcolutionClient) Create(ctx context.Context, tree Tree) error {
	body := client.TreeCreate{
		Description:  fmt.Sprintf("%s - Dieser Baum wurde importiert", tree.Description),
		Latitude:     tree.Latitude,
		Longitude:    tree.Longitude,
		Number:       tree.Number,
		PlantingYear: tree.PlantingYear,
		Species:      tree.Species,
		Provider:     &r.provider,
	}

	_, _, err := r.client.TreeAPI.CreateTree(ctx).Body(body).Execute()
	return err
}

func (r *GreenEcolutionClient) Update(ctx context.Context, id int32, tree Tree) error {
	body := client.TreeUpdate{
		Description:  tree.Description,
		Latitude:     tree.Latitude,
		Longitude:    tree.Longitude,
		Number:       tree.Number,
		PlantingYear: tree.PlantingYear,
		Species:      tree.Species,
		Provider:     &r.provider,
	}

	_, _, err := r.client.TreeAPI.UpdateTree(ctx, id).Body(body).Execute()
	return err
}

func (r *GreenEcolutionClient) Delete(ctx context.Context, id int32) error {
	_, err := r.client.TreeAPI.DeleteTree(ctx, id).Execute()
	return err
}

func (r *GreenEcolutionClient) Archive(ctx context.Context, id int32) error {
	_, err := r.client.TreeAPI.DeleteTree(ctx, id).Execute() // TODO: Archive
	return err
}

func (r *GreenEcolutionClient) RefreshToken(ctx context.Context, refreshToken string) (*client.ClientToken, error) {
	body := client.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	token, _, err := r.client.UserAPI.V1UserTokenRefreshPost(ctx).Body(body).Execute()
	return token, err
}
