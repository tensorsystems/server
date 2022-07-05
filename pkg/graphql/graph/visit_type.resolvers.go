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

func (r *mutationResolver) SaveVisitType(ctx context.Context, input graph_models.VisitTypeInput) (*models.VisitType, error) {
	var visitType repository.VisitType
	deepCopy.Copy(&input).To(&visitType)

	err := visitType.Save()
	if err != nil {
		return nil, err
	}

	return &visitType, nil
}

func (r *mutationResolver) UpdateVisitType(ctx context.Context, input graph_models.VisitTypeInput, id int) (*models.VisitType, error) {
	var visitType repository.VisitType
	deepCopy.Copy(&input).To(&visitType)
	visitType.ID = id

	if err := visitType.Update(); err != nil {
		return nil, err
	}

	return &visitType, nil
}

func (r *mutationResolver) DeleteVisitType(ctx context.Context, id int) (bool, error) {
	var visitType repository.VisitType
	err := visitType.Delete(id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) VisitTypes(ctx context.Context, page models.PaginationInput) (*graph_models.VisitTypeConnection, error) {
	var visitType repository.VisitType
	visitTypes, count, err := visitType.GetAll(page)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.VisitTypeEdge, len(visitTypes))

	for i, entity := range visitTypes {
		e := entity

		edges[i] = &model.VisitTypeEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(visitTypes, count, page)
	return &model.VisitTypeConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
