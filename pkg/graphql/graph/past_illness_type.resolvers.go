package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SavePastIllnessTypes(ctx context.Context, input model.PastIllnessTypeInput) (*repository.PastIllnessType, error) {
	var entity repository.PastIllnessType
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdatePastIllnessType(ctx context.Context, input model.PastIllnessTypeUpdateInput) (*repository.PastIllnessType, error) {
	var entity repository.PastIllnessType
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeletePastIllnessType(ctx context.Context, id int) (bool, error) {
	var entity repository.PastIllnessType

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) PastIllnessType(ctx context.Context, id int) (*repository.PastIllnessType, error) {
	var entity repository.PastIllnessType

	if err := entity.Get(id); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) PastIllnessTypes(ctx context.Context, page repository.PaginationInput) (*model.PastIllnessTypeConnection, error) {
	var entity repository.PastIllnessType
	result, count, err := entity.GetAll(page)
	if err != nil {
		return nil, err
	}

	edges := make([]*model.PastIllnessTypeEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &model.PastIllnessTypeEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &model.PastIllnessTypeConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
