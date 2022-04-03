package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/Tensor-Systems/tensoremr-server/pkg/graphql/graph/model"
	"github.com/Tensor-Systems/tensoremr-server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) CreateEyewearShop(ctx context.Context, input model.EyewearShopInput) (*repository.EyewearShop, error) {
	var entity repository.EyewearShop
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateEyewearShop(ctx context.Context, input model.EyewearShopUpdateInput) (*repository.EyewearShop, error) {
	var entity repository.EyewearShop
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteEyewearShop(ctx context.Context, id int) (bool, error) {
	var entity repository.EyewearShop
	if err := entity.Delete(id); err != nil {
		return false, err
	}
	return true, nil
}

func (r *queryResolver) EyewearShop(ctx context.Context, id int) (*repository.EyewearShop, error) {
	var entity repository.EyewearShop
	if err := entity.Get(id); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *queryResolver) EyewearShops(ctx context.Context, page repository.PaginationInput) (*model.EyewearShopConnection, error) {
	var entity repository.EyewearShop
	result, count, err := entity.GetAll(page, nil)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.EyewearShopEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &model.EyewearShopEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &model.EyewearShopConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
