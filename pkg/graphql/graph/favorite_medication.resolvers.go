package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/middleware"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SaveFavoriteMedication(ctx context.Context, input model.FavoriteMedicationInput) (*repository.FavoriteMedication, error) {
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

	var entity repository.FavoriteMedication
	deepCopy.Copy(&input).To(&entity)

	entity.UserID = user.ID

	if err := entity.Save(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateFavoriteMedication(ctx context.Context, input model.FavoriteMedicationUpdateInput) (*repository.FavoriteMedication, error) {
	var entity repository.FavoriteMedication
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteFavoriteMedication(ctx context.Context, id int) (bool, error) {
	var entity repository.FavoriteMedication

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) FavoriteMedications(ctx context.Context, page repository.PaginationInput, filter *model.FavoriteMedicationFilter, searchTerm *string) (*model.FavoriteMedicationConnection, error) {
	var f repository.FavoriteMedication
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	var entity repository.FavoriteMedication
	entities, count, err := entity.GetAll(page, &f, searchTerm)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.FavoriteMedicationEdge, len(entities))

	for i, entity := range entities {
		e := entity

		edges[i] = &model.FavoriteMedicationEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(entities, count, page)
	return &model.FavoriteMedicationConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) UserFavoriteMedications(ctx context.Context, page repository.PaginationInput, searchTerm *string) (*model.FavoriteMedicationConnection, error) {
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

	var entity repository.FavoriteMedication
	entities, count, err := entity.GetAll(page, &repository.FavoriteMedication{UserID: user.ID}, searchTerm)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.FavoriteMedicationEdge, len(entities))

	for i, entity := range entities {
		e := entity

		edges[i] = &model.FavoriteMedicationEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(entities, count, page)
	return &model.FavoriteMedicationConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) SearchFavoriteMedications(ctx context.Context, searchTerm string, page repository.PaginationInput) (*model.FavoriteMedicationConnection, error) {
	var entity repository.FavoriteMedication
	entities, count, err := entity.Search(page, searchTerm)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.FavoriteMedicationEdge, len(entities))

	for i, entity := range entities {
		e := entity

		edges[i] = &model.FavoriteMedicationEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(entities, count, page)
	return &model.FavoriteMedicationConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
