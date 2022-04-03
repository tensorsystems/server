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

func (r *mutationResolver) SaveChiefComplaintType(ctx context.Context, input model.ChiefComplaintTypeInput) (*repository.ChiefComplaintType, error) {
	var entity repository.ChiefComplaintType
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateChiefComplaintType(ctx context.Context, input model.ChiefComplaintTypeUpdateInput) (*repository.ChiefComplaintType, error) {
	var entity repository.ChiefComplaintType
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteChiefComplaintType(ctx context.Context, id int) (bool, error) {
	var entity repository.ChiefComplaintType

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) ChiefComplaintType(ctx context.Context, id int) (*repository.ChiefComplaintType, error) {
	var entity repository.ChiefComplaintType

	if err := entity.Get(id); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) ChiefComplaintTypes(ctx context.Context, page repository.PaginationInput, searchTerm *string, favorites *bool) (*model.ChiefComplaintTypeConnection, error) {
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

	var result []repository.ChiefComplaintType
	var count int64

	var entity repository.ChiefComplaintType

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

	edges := make([]*model.ChiefComplaintTypeEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &model.ChiefComplaintTypeEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &model.ChiefComplaintTypeConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
