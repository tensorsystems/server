package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SaveSupply(ctx context.Context, input model.SupplyInput) (*repository.Supply, error) {
	var entity repository.Supply
	deepCopy.Copy(&input).To(&entity)

	var billing repository.Billing
	billings, err := billing.GetByIds(input.BillingIds)
	if err != nil {
		return nil, err
	}

	entity.Billings = billings

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateSupply(ctx context.Context, input model.SupplyUpdateInput) (*repository.Supply, error) {
	var entity repository.Supply
	deepCopy.Copy(&input).To(&entity)

	var billing repository.Billing
	billings, err := billing.GetByIds(input.BillingIds)
	if err != nil {
		return nil, err
	}

	entity.Billings = billings

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteSupply(ctx context.Context, id int) (bool, error) {
	var entity repository.Supply

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) Supplies(ctx context.Context, page repository.PaginationInput, searchTerm *string) (*model.SupplyConnection, error) {
	var entity repository.Supply
	result, count, err := entity.GetAll(page, searchTerm)
	if err != nil {
		return nil, err
	}

	edges := make([]*model.SupplyEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &model.SupplyEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &model.SupplyConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
