package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	graph_models "github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/models"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SaveHpiComponentType(ctx context.Context, input graph_models.HpiComponentTypeInput) (*models.HpiComponentType, error) {
	var hpiComponentType models.HpiComponentType
	deepCopy.Copy(&input).To(&hpiComponentType)

	var repository repository.HpiComponentTypeRepository
	if err := repository.Save(&hpiComponentType); err != nil {
		return nil, err
	}

	return &hpiComponentType, nil
}

func (r *mutationResolver) UpdateHpiComponentType(ctx context.Context, input graph_models.HpiComponentTypeUpdateInput) (*models.HpiComponentType, error) {
	var hpiComponentType models.HpiComponentType
	deepCopy.Copy(&input).To(&hpiComponentType)

	var repository repository.HpiComponentTypeRepository
	if err := repository.Update(&hpiComponentType); err != nil {
		return nil, err
	}

	return &hpiComponentType, nil
}

func (r *mutationResolver) DeleteHpiComponentType(ctx context.Context, id int) (bool, error) {
	var repository repository.HpiComponentTypeRepository
	
	if err := repository.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) SaveHpiComponent(ctx context.Context, input graph_models.HpiComponentInput) (*models.HpiComponent, error) {
	var hpiComponent models.HpiComponent
	deepCopy.Copy(&input).To(&hpiComponent)

	var repository repository.HpiComponentRepository
	if err := repository.Save(&hpiComponent); err != nil {
		return nil, err
	}

	return &hpiComponent, nil
}

func (r *mutationResolver) UpdateHpiComponent(ctx context.Context, input graph_models.HpiComponentUpdateInput) (*models.HpiComponent, error) {
	var hpiComponent models.HpiComponent
	deepCopy.Copy(&input).To(&hpiComponent)

	var repository repository.HpiComponentRepository
	if err := repository.Update(&hpiComponent); err != nil {
		return nil, err
	}

	return &hpiComponent, nil
}

func (r *mutationResolver) DeleteHpiComponent(ctx context.Context, id int) (bool, error) {
	var repository repository.HpiComponentRepository

	if err := repository.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) HpiComponentTypes(ctx context.Context, page models.PaginationInput) (*graph_models.HpiComponentTypeConnection, error) {
	var repository repository.HpiComponentTypeRepository

	hpiComponentTypes, count, err := repository.GetAll(page)

	if err != nil {
		return nil, err
	}

	edges := make([]*graph_models.HpiComponentTypeEdge, len(hpiComponentTypes))

	for i, entity := range hpiComponentTypes {
		e := entity

		edges[i] = &graph_models.HpiComponentTypeEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(hpiComponentTypes, count, page)
	return &graph_models.HpiComponentTypeConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) HpiComponents(ctx context.Context, page models.PaginationInput, filter *graph_models.HpiFilter, searchTerm *string) (*graph_models.HpiComponentConnection, error) {
	var f models.HpiComponent
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	var repository repository.HpiComponentRepository

	hpiComponents, count, err := repository.Search(page, &f, searchTerm)

	if err != nil {
		return nil, err
	}

	edges := make([]*graph_models.HpiComponentEdge, len(hpiComponents))

	for i, entity := range hpiComponents {
		e := entity

		edges[i] = &graph_models.HpiComponentEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(hpiComponents, count, page)
	return &graph_models.HpiComponentConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
