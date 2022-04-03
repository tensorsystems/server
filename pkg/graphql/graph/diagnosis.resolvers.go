package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/Tensor-Systems/tensoremr-server/pkg/graphql/graph/model"
	"github.com/Tensor-Systems/tensoremr-server/pkg/middleware"
	"github.com/Tensor-Systems/tensoremr-server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SaveDiagnosis(ctx context.Context, input model.DiagnosisInput) (*repository.Diagnosis, error) {
	var entity repository.Diagnosis
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateDiagnosis(ctx context.Context, input model.DiagnosisUpdateInput) (*repository.Diagnosis, error) {
	var entity repository.Diagnosis
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteDiagnosis(ctx context.Context, id int) (bool, error) {
	var entity repository.Diagnosis

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) Diagnoses(ctx context.Context, page repository.PaginationInput, searchTerm *string, favorites *bool) (*model.DiagnosisConnection, error) {
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return nil, errors.New("Cannot find user")
	}

	var user repository.User
	err = user.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	var entity repository.Diagnosis

	var result []repository.Diagnosis
	var count int64

	if favorites != nil && *favorites == true {
		result, count, err = entity.GetFavorites(page, searchTerm, user.ID)
		if err != nil {
			return nil, err
		}
	} else {
		result, count, err = entity.GetAll(page, searchTerm)
		if err != nil {
			return nil, err
		}
	}

	edges := make([]*model.DiagnosisEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &model.DiagnosisEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &model.DiagnosisConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
