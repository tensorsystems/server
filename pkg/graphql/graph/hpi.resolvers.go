package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SaveHpiComponentType(ctx context.Context, input model.HpiComponentTypeInput) (*repository.HpiComponentType, error) {
	var hpiComponentType repository.HpiComponentType

	deepCopy.Copy(&input).To(&hpiComponentType)

	err := hpiComponentType.Save()
	if err != nil {
		return nil, err
	}

	return &hpiComponentType, nil
}

func (r *mutationResolver) UpdateHpiComponentType(ctx context.Context, input model.HpiComponentTypeUpdateInput) (*repository.HpiComponentType, error) {
	var hpiComponentType repository.HpiComponentType

	deepCopy.Copy(&input).To(&hpiComponentType)

	if err := hpiComponentType.Update(); err != nil {
		return nil, err
	}

	return &hpiComponentType, nil
}

func (r *mutationResolver) DeleteHpiComponentType(ctx context.Context, id int) (bool, error) {
	var hpiComponentType repository.HpiComponentType

	err := hpiComponentType.Delete(id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) SaveHpiComponent(ctx context.Context, input model.HpiComponentInput) (*repository.HpiComponent, error) {
	var hpiComponent repository.HpiComponent

	deepCopy.Copy(&input).To(&hpiComponent)

	err := hpiComponent.Save()
	if err != nil {
		return nil, err
	}

	return &hpiComponent, nil
}

func (r *mutationResolver) UpdateHpiComponent(ctx context.Context, input model.HpiComponentUpdateInput) (*repository.HpiComponent, error) {
	var hpiComponent repository.HpiComponent

	deepCopy.Copy(&input).To(&hpiComponent)

	if err := hpiComponent.Update(); err != nil {
		return nil, err
	}

	return &hpiComponent, nil
}

func (r *mutationResolver) DeleteHpiComponent(ctx context.Context, id int) (bool, error) {
	var hpiComponent repository.HpiComponent

	err := hpiComponent.Delete(id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) HpiComponentTypes(ctx context.Context, page repository.PaginationInput) (*model.HpiComponentTypeConnection, error) {
	var hpiComponentType repository.HpiComponentType

	hpiComponentTypes, count, err := hpiComponentType.GetAll(page)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.HpiComponentTypeEdge, len(hpiComponentTypes))

	for i, entity := range hpiComponentTypes {
		e := entity

		edges[i] = &model.HpiComponentTypeEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(hpiComponentTypes, count, page)
	return &model.HpiComponentTypeConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) HpiComponents(ctx context.Context, page repository.PaginationInput, filter *model.HpiFilter, searchTerm *string) (*model.HpiComponentConnection, error) {
	var f repository.HpiComponent

	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	hpiComponents, count, err := f.Search(page, &f, searchTerm)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.HpiComponentEdge, len(hpiComponents))

	for i, entity := range hpiComponents {
		e := entity

		edges[i] = &model.HpiComponentEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(hpiComponents, count, page)
	return &model.HpiComponentConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
