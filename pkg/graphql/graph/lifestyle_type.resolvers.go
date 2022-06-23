package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SaveLifestyleTypes(ctx context.Context, input model.LifestyleTypeInput) (*repository.LifestyleType, error) {
	var entity repository.LifestyleType
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateLifestyleType(ctx context.Context, input model.LifestyleTypeUpdateInput) (*repository.LifestyleType, error) {
	var entity repository.LifestyleType
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteLifestyleType(ctx context.Context, id int) (bool, error) {
	var entity repository.LifestyleType

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) LifestyleType(ctx context.Context, id int) (*repository.LifestyleType, error) {
	var entity repository.LifestyleType

	if err := entity.Get(id); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) LifestyleTypes(ctx context.Context, page repository.PaginationInput) (*model.LifestyleTypeConnection, error) {
	var entity repository.LifestyleType
	result, count, err := entity.GetAll(page)
	if err != nil {
		return nil, err
	}

	edges := make([]*model.LifestyleTypeEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &model.LifestyleTypeEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &model.LifestyleTypeConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
