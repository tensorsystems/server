package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SavePatientEncounterLimit(ctx context.Context, input model.PatientEncounterLimitInput) (*repository.PatientEncounterLimit, error) {
	var entity repository.PatientEncounterLimit
	deepCopy.Copy(&input).To(&entity)

	entity.Overbook = 5

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdatePatientEncounterLimit(ctx context.Context, input model.PatientEncounterLimitUpdateInput) (*repository.PatientEncounterLimit, error) {
	var entity repository.PatientEncounterLimit
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeletePatientEncounterLimit(ctx context.Context, id int) (bool, error) {
	var entity repository.PatientEncounterLimit

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) PatientEncounterLimit(ctx context.Context, id int) (*repository.PatientEncounterLimit, error) {
	var entity repository.PatientEncounterLimit

	if err := entity.Get(id); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) PatientEncounterLimits(ctx context.Context, page repository.PaginationInput) (*model.PatientEncounterLimitConnection, error) {
	var entity repository.PatientEncounterLimit
	result, count, err := entity.GetAll(page)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.PatientEncounterLimitEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &model.PatientEncounterLimitEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &model.PatientEncounterLimitConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) PatientEncounterLimitByUser(ctx context.Context, userID int) (*repository.PatientEncounterLimit, error) {
	var entity repository.PatientEncounterLimit

	if err := entity.GetByUser(userID); err != nil {
		return nil, err
	}

	return &entity, nil
}
