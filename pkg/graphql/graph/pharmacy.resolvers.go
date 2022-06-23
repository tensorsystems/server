package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) CreatePharmacy(ctx context.Context, input model.PharmacyInput) (*repository.Pharmacy, error) {
	var entity repository.Pharmacy
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdatePharmacy(ctx context.Context, input model.PharmacyUpdateInput) (*repository.Pharmacy, error) {
	var entity repository.Pharmacy
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeletePharmacy(ctx context.Context, id int) (bool, error) {
	var entity repository.Pharmacy
	if err := entity.Delete(id); err != nil {
		return false, err
	}
	return true, nil
}

func (r *queryResolver) Pharmacy(ctx context.Context, id int) (*repository.Pharmacy, error) {
	var entity repository.Pharmacy
	if err := entity.Get(id); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *queryResolver) Pharmacies(ctx context.Context, page repository.PaginationInput) (*model.PharmacyConnection, error) {
	var entity repository.Pharmacy
	result, count, err := entity.GetAll(page, nil)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.PharmacyEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &model.PharmacyEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &model.PharmacyConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
